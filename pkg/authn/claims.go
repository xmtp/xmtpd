package authn

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// XMTPD_COMPATIBLE_VERSION_CONSTRAINT major or minor version bumps indicate backwards incompatible changes
	XMTPD_COMPATIBLE_VERSION_CONSTRAINT = "~ 0.1.3"
)

type XmtpdClaims struct {
	Version *string `json:"version,omitempty"`
	jwt.RegisteredClaims
}

func ValidateVersionClaimIsCompatible(claims *XmtpdClaims) error {
	if claims.Version == nil || *claims.Version == "" {
		return nil
	}

	c, err := semver.NewConstraint(XMTPD_COMPATIBLE_VERSION_CONSTRAINT)
	if err != nil {
		return err
	}

	v, err := semver.NewVersion(*claims.Version)

	if err != nil {
		return err
	}

	if ok := c.Check(v); !ok {
		return fmt.Errorf("version %s is not compatible", *claims.Version)
	}

	return nil
}
