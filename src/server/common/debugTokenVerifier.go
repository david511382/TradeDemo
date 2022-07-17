package common

import (
	"zerologix-homework/src/server/domain"
)

type debugTokenVerifier struct {
	json jsonTokenVerifier
}

func NewDebugTokenVerifier() debugTokenVerifier {
	return debugTokenVerifier{
		json: NewJsonTokenVerifier(),
	}
}

// token is json of JwtClaims
func (l debugTokenVerifier) Parse(token string) (jwtClaims domain.JwtClaims, resultErr error) {
	claims, err := l.json.Parse(token)
	jwtClaims = claims
	if err != nil {
		resultErr = err
	}

	return
}
