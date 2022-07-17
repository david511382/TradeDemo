package common

import (
	"encoding/json"
	errUtil "zerologix-homework/src/pkg/util/error"
	"zerologix-homework/src/server/domain"
)

type jsonTokenVerifier struct {
}

func NewJsonTokenVerifier() jsonTokenVerifier {
	return jsonTokenVerifier{}
}

// token is json of JwtClaims
func (l jsonTokenVerifier) Parse(token string) (jwtClaims domain.JwtClaims, resultErr error) {
	jwtClaimsP := &domain.JwtClaims{}
	if err := json.Unmarshal([]byte(token), jwtClaimsP); err != nil {
		err := errUtil.NewError(err)
		resultErr = err
		return
	}
	jwtClaims = *jwtClaimsP
	return
}
