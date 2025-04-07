package authn

type CloseFunc func()

type JWTVerifier interface {
	Verify(tokenString string) (uint32, CloseFunc, error)
}
