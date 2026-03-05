package authn

import (
	"errors"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/metrics"

	"github.com/Masterminds/semver/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// MinCompatibleVersion is the minimum peer version this node will accept connections from.
// Update this constant when introducing breaking wire protocol changes; do not change it
// for routine minor releases.
var MinCompatibleVersion = semver.MustParse("1.1.0")

type XmtpdClaims struct {
	Version *semver.Version `json:"version,omitempty"`
	jwt.RegisteredClaims
}

type ClaimValidator struct {
	constraint semver.Constraints
}

func NewClaimValidator(logger *zap.Logger, serverVersion *semver.Version) (*ClaimValidator, error) {
	if serverVersion == nil {
		return nil, errors.New("serverVersion is nil")
	}

	// Accept peers in [floor, next major) where floor is MinCompatibleVersion when the server
	// is on the same major, or major.0.0 otherwise. Minor bumps are backward-compatible;
	// only a major bump or an explicit MinCompatibleVersion update signals a breaking change.
	floor := MinCompatibleVersion.String()
	if serverVersion.Major() != MinCompatibleVersion.Major() {
		floor = fmt.Sprintf("%d.0.0", serverVersion.Major())
	}
	constraintStr := fmt.Sprintf(">=%s, <%d.0.0", floor, serverVersion.Major()+1)

	logger.Debug(
		"using semver constraint for sync compatibility",
		zap.String("constraint", constraintStr),
	)

	constraint, err := semver.NewConstraint(constraintStr)
	if err != nil {
		return nil, err
	}

	return &ClaimValidator{constraint: *constraint}, nil
}

func (cv *ClaimValidator) ValidateVersionClaimIsCompatible(claims *XmtpdClaims) (CloseFunc, error) {
	if claims.Version == nil {
		return emptyClose, nil
	}

	// SemVer implementations generally do not consider pre-releases to be valid next releases
	// we use SemVer here to allow incoming connections, for which in-development nodes are acceptable
	// see discussion in https://github.com/Masterminds/semver/issues/21
	sanitizedVersion, err := claims.Version.SetPrerelease("")
	if err != nil {
		return emptyClose, err
	}

	metrics.EmitNewConnectionRequestVersion(sanitizedVersion.String())

	if ok := cv.constraint.Check(&sanitizedVersion); !ok {
		return emptyClose, fmt.Errorf("serverVersion %s is not compatible", *claims.Version)
	}

	tracker := metrics.NewIncomingConnectionTracker(sanitizedVersion.String())
	tracker.Open()

	return func() {
		tracker.Close()
	}, nil
}
