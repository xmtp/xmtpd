// Package authn implements a JWT token factory and verifier.
package authn

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmtp/xmtpd/pkg/registry"
)

const (
	maxTokenDuration = 2 * time.Hour
	MaxClockSkew     = 2 * time.Minute
)

type RegistryVerifier struct {
	registry  registry.NodeRegistry
	myNodeID  uint32
	validator ClaimValidator
}

func emptyClose() {}

// NewRegistryVerifier returns a new RegistryVerifier that connects to the NodeRegistry and verifies JWTs
// against the registered public keys based on the JWT's subject field.
func NewRegistryVerifier(
	logger *zap.Logger,
	registry registry.NodeRegistry,
	myNodeID uint32,
	serverVersion *semver.Version,
) (*RegistryVerifier, error) {
	validator, err := NewClaimValidator(logger, serverVersion)
	if err != nil {
		return nil, err
	}

	return &RegistryVerifier{registry: registry, myNodeID: myNodeID, validator: *validator}, nil
}

func (v *RegistryVerifier) Verify(tokenString string) (uint32, CloseFunc, error) {
	var token *jwt.Token
	var err error

	if token, err = jwt.ParseWithClaims(
		tokenString,
		&XmtpdClaims{},
		v.getMatchingPublicKey,
	); err != nil {
		return 0, emptyClose, err
	}
	if err = v.validateAudience(token); err != nil {
		return 0, emptyClose, err
	}

	if err = validateExpiry(token); err != nil {
		return 0, emptyClose, err
	}

	closer, err := v.validateClaims(token)
	if err != nil {
		return 0, emptyClose, err
	}

	nodeID, err := getSubjectNodeID(token)
	if err != nil {
		return 0, emptyClose, err
	}

	return nodeID, closer, nil
}

func (v *RegistryVerifier) getMatchingPublicKey(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*SigningMethodSecp256k1); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	subjectNodeID, err := getSubjectNodeID(token)
	if err != nil {
		return nil, fmt.Errorf("could not get subject node ID: %w", err)
	}

	node, err := v.registry.GetNode(subjectNodeID)
	if err != nil {
		return nil, fmt.Errorf("could not get node %d from registry: %w", subjectNodeID, err)
	}

	return node.SigningKey, nil
}

// Ensure that the token was intended for this node and not replayed from another node
// The audience field must be this node's ID
func (v *RegistryVerifier) validateAudience(token *jwt.Token) error {
	audience, err := token.Claims.GetAudience()
	if err != nil {
		return err
	}

	for _, audienceString := range audience {
		audienceNodeID, err := parseInt32(audienceString)
		if err != nil {
			continue
		}

		if audienceNodeID == v.myNodeID {
			return nil
		}
	}

	return fmt.Errorf("could not find node ID in audience %v", audience)
}

func (v *RegistryVerifier) validateClaims(token *jwt.Token) (CloseFunc, error) {
	claims, ok := token.Claims.(*XmtpdClaims)
	if !ok {
		return emptyClose, errors.New("invalid token claims type")
	}

	// Check if the token is valid
	if !token.Valid {
		return emptyClose, errors.New("invalid token")
	}

	return v.validator.ValidateVersionClaimIsCompatible(claims)
}

// Parse the subject claim of the JWT and return the node ID as a uint32
func getSubjectNodeID(token *jwt.Token) (uint32, error) {
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	nodeID, err := parseInt32(subject)
	if err != nil {
		return 0, err
	}

	return nodeID, nil
}

// Validate the issued at and expiration time claims of the JWT
func validateExpiry(token *jwt.Token) error {
	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return err
	}
	issuedAt, err := token.Claims.GetIssuedAt()
	if err != nil {
		return err
	}

	// We allow tokens to be issued up to 2 minutes in the future to account for clock skew
	if time.Since(issuedAt.Time) < MaxClockSkew*-1 {
		return errors.New("token issued in the future")
	}

	// Tokens cannot expire before they are issued
	if exp.Before(issuedAt.Time) {
		return errors.New("token expires before the issued at time")
	}

	// Tokens can only have a validity period of at most 2 hours
	if exp.Sub(issuedAt.Time) > maxTokenDuration {
		return errors.New("token expiration time is greater than the max duration")
	}

	// Tokens cannot be expired
	if time.Since(exp.Time) > 0 {
		return errors.New("token is expired")
	}

	return nil
}

func parseInt32(str string) (uint32, error) {
	parsed, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(parsed), nil
}
