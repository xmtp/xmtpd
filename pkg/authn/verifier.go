package authn

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"go.uber.org/zap"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmtp/xmtpd/pkg/registry"
)

const (
	MAX_TOKEN_DURATION = 2 * time.Hour
	MAX_CLOCK_SKEW     = 2 * time.Minute
)

type RegistryVerifier struct {
	registry  registry.NodeRegistry
	myNodeID  uint32
	validator ClaimValidator
}

/*
A RegistryVerifier connects to the NodeRegistry and verifies JWTs against the registered public keys
based on the JWT's subject field
*/
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

func (v *RegistryVerifier) Verify(tokenString string) (uint32, error) {
	var token *jwt.Token
	var err error

	if token, err = jwt.ParseWithClaims(
		tokenString,
		&XmtpdClaims{},
		v.getMatchingPublicKey,
	); err != nil {
		return 0, err
	}
	if err = v.validateAudience(token); err != nil {
		return 0, err
	}

	if err = validateExpiry(token); err != nil {
		return 0, err
	}

	if err = v.validateClaims(token); err != nil {
		return 0, err
	}

	return getSubjectNodeId(token)
}

func (v *RegistryVerifier) getMatchingPublicKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*SigningMethodSecp256k1); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	subjectNodeId, err := getSubjectNodeId(token)
	if err != nil {
		return nil, err
	}

	node, err := v.registry.GetNode(subjectNodeId)
	if err != nil {
		return nil, err
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
		audienceNodeId, err := parseInt32(audienceString)
		if err != nil {
			continue
		}

		if audienceNodeId == v.myNodeID {
			return nil
		}
	}

	return fmt.Errorf("could not find node ID in audience %v", audience)
}

func (v *RegistryVerifier) validateClaims(token *jwt.Token) error {
	claims, ok := token.Claims.(*XmtpdClaims)
	if !ok {
		return fmt.Errorf("invalid token claims type")
	}

	// Check if the token is valid
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return v.validator.ValidateVersionClaimIsCompatible(claims)
}

// Parse the subject claim of the JWT and return the node ID as a uint32
func getSubjectNodeId(token *jwt.Token) (uint32, error) {
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	nodeId, err := parseInt32(subject)
	if err != nil {
		return 0, err
	}

	return uint32(nodeId), nil
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
	if time.Since(issuedAt.Time) < MAX_CLOCK_SKEW*-1 {
		return fmt.Errorf("token issued in the future")
	}

	// Tokens cannot expire before they are issued
	if exp.Before(issuedAt.Time) {
		return fmt.Errorf("token expires before the issued at time")
	}

	// Tokens can only have a validity period of at most 2 hours
	if exp.Sub(issuedAt.Time) > MAX_TOKEN_DURATION {
		return fmt.Errorf("token expiration time is greater than the max duration")
	}

	// Tokens cannot be expired
	if time.Since(exp.Time) > 0 {
		return fmt.Errorf("token is expired")
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
