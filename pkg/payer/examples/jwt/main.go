package main

import (
	"context"
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmtp/xmtpd/pkg/payer"
)

const EXPECTED_ISSUER = "my-app.com"

var (
	ErrMissingToken     = errors.New("missing JWT token")
	ErrInvalidToken     = errors.New("invalid JWT token")
	ErrInvalidSignature = errors.New("invalid token signature")
)

// jwtIdentityFn creates an identity function that verifies JWTs
func jwtIdentityFn(publicKey []byte) payer.IdentityFn {
	return func(ctx context.Context) (payer.Identity, error) {
		authHeader := payer.AuthorizationHeaderFromContext(ctx)
		if authHeader == "" {
			return payer.Identity{}, payer.NewUnauthenticatedError(
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
					return nil, payer.NewPermissionDeniedError(
						"Invalid signing method",
						ErrInvalidSignature,
					)
				}
				return publicKey, nil
			},
			jwt.WithIssuer(EXPECTED_ISSUER),
		)
		if err != nil {
			return payer.Identity{}, payer.NewPermissionDeniedError("failed to validate token", err)
		}

		// Extract claims
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			return payer.Identity{}, payer.NewPermissionDeniedError(
				"failed to validate token",
				ErrInvalidToken,
			)
		}

		userID, err := claims.GetSubject()
		if err != nil {
			return payer.Identity{}, payer.NewPermissionDeniedError(
				"failed to get subject from token",
				err,
			)
		}

		// Return identity based on JWT claims
		return payer.NewUserIdentity(userID), nil
	}
}

func main() {
	// In a real application, this would be a secure key loaded from environment/config
	publicKey := []byte("your-applications-public-key")

	payerService, err := payer.NewPayerServiceBuilder(payer.MustLoadConfig()).
		WithIdentityFn(jwtIdentityFn(publicKey)).
		Build()
	if err != nil {
		log.Fatalf("Failed to build payer service: %v", err)
	}

	payerService.WaitForShutdown()
}
