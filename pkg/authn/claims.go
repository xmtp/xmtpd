package authn

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type XmtpdClaims struct {
	Version *semver.Version `json:"version,omitempty"`
	jwt.RegisteredClaims
}
type ClaimValidator struct {
	constraint semver.Constraints
}

func NewClaimValidator(logger *zap.Logger, serverVersion *semver.Version) (*ClaimValidator, error) {
	if serverVersion == nil {
		return nil, fmt.Errorf("serverVersion is nil")
	}

	// https://github.com/Masterminds/semver?tab=readme-ov-file#caret-range-comparisons-major
	constraintStr := fmt.Sprintf("^%d.%d", serverVersion.Major(), serverVersion.Minor())
	logger.Debug(
		"Using semver constraint for sync compatibility",
		zap.String("constraint", constraintStr),
	)

	constraint, err := semver.NewConstraint(constraintStr)
	if err != nil {
		return nil, err
	}

	return &ClaimValidator{constraint: *constraint}, nil
}
func (cv *ClaimValidator) ValidateVersionClaimIsCompatible(claims *XmtpdClaims) error {
	if claims.Version == nil {
		return nil
	}

	// SemVer implementations generally do not consider pre-releases to be valid next releases
	// we use SemVer here to allow incoming connections, for which in-development nodes are acceptable
	// see discussion in https://github.com/Masterminds/semver/issues/21
	sanitizedVersion, err := claims.Version.SetPrerelease("")
	if err != nil {
		return err
	}
	if ok := cv.constraint.Check(&sanitizedVersion); !ok {
		return fmt.Errorf("serverVersion %s is not compatible", *claims.Version)
	}

	return nil
}
