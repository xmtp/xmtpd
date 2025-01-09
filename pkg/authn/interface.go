package authn

type JWTVerifier interface {
	Verify(tokenString string) (uint32, error)
}
