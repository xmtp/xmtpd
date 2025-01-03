package authn

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// XMTPD_COMPATIBLE_VERSION_CONSTRAINT major or minor serverVersion bumps indicate backwards incompatible changes
	XMTPD_COMPATIBLE_VERSION_CONSTRAINT = "~ 0.1.3"
)

type XmtpdClaims struct {
	Version *semver.Version `json:"version,omitempty"`
	jwt.RegisteredClaims
}

func ValidateVersionClaimIsCompatible(claims *XmtpdClaims) error {
	if claims.Version == nil {
		return nil
	}

	c, err := semver.NewConstraint(XMTPD_COMPATIBLE_VERSION_CONSTRAINT)
	if err != nil {
		return err
	}

	// SemVer implementations generally do not consider pre-releases to be valid next releases
	// we use SemVer here to allow incoming connections, for which in-development nodes are acceptable
	// see discussion in https://github.com/Masterminds/semver/issues/21
	sanitizedVersion, err := claims.Version.SetPrerelease("")
	if err != nil {
		return err
	}
	if ok := c.Check(&sanitizedVersion); !ok {
		return fmt.Errorf("serverVersion %s is not compatible", *claims.Version)
	}

	return nil
}
