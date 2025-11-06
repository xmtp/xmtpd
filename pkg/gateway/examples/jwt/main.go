package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmtp/xmtpd/pkg/gateway"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const expectedIssuer = "my-app.com"

var (
	ErrMissingToken     = errors.New("missing JWT token")
	ErrInvalidToken     = errors.New("invalid JWT token")
	ErrInvalidSignature = errors.New("invalid token signature")
)

// jwtIdentityFn creates an identity function that verifies JWTs
func jwtIdentityFn(publicKey []byte) gateway.IdentityFn {
	return func(headers http.Header, _ string) (gateway.Identity, error) {
		authHeader := utils.AuthorizationTokenFromHeader(headers)
		if authHeader == "" {
			return gateway.Identity{}, gateway.NewUnauthenticatedError(
				"Missing JWT token",
				ErrMissingToken,
			)
		}

		// Parse and verify the token
		token, err := jwt.ParseWithClaims(
			authHeader,
			&jwt.RegisteredClaims{},
			func(token *jwt.Token) (interface{}, error) {
				// Verify signing method
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, gateway.NewPermissionDeniedError(
						"Invalid signing method",
						ErrInvalidSignature,
					)
				}
				return publicKey, nil
			},
			jwt.WithIssuer(expectedIssuer),
		)
		if err != nil {
			return gateway.Identity{}, gateway.NewPermissionDeniedError(
				"failed to validate token",
				err,
			)
		}

		// Extract claims
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			return gateway.Identity{}, gateway.NewPermissionDeniedError(
				"failed to validate token",
				ErrInvalidToken,
			)
		}

		userID, err := claims.GetSubject()
		if err != nil {
			return gateway.Identity{}, gateway.NewPermissionDeniedError(
				"failed to get subject from token",
				err,
			)
		}

		// Return identity based on JWT claims
		return gateway.NewUserIdentity(userID), nil
	}
}

func main() {
	// In a real application, this would be a secure key loaded from environment/config
	publicKey := []byte("your-applications-public-key")

	gatewayService, err := gateway.NewGatewayServiceBuilder(gateway.MustLoadConfig()).
		WithIdentityFn(jwtIdentityFn(publicKey)).
		Build()
	if err != nil {
		log.Fatalf("failed to build gateway service: %v", err)
	}

	gatewayService.WaitForShutdown(30 * time.Second)
}
