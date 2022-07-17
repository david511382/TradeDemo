package domain

type ITokenVerifier interface {
	Parse(token string) (jwtClaims JwtClaims, resultErr error)
}
