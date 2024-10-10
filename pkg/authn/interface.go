package authn

type JWTVerifier interface {
	Verify(tokenString string) error
}
