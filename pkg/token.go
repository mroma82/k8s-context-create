package pkg

import (
	"fmt"

	"github.com/lestrrat-go/jwx/jwt"
)

// type
type Token struct {
	Val            string
	Namespace      string
	ServiceAccount string
	Secret         string
}

// function that parses a token
func ParseToken(tokenStr string) (*Token, error) {

	// get the token
	t, err := jwt.Parse([]byte(tokenStr))
	if err != nil {
		return nil, err
	}

	// init token
	token := &Token{
		Val: tokenStr,
	}

	// namespace
	if v, ok := t.Get("kubernetes.io/serviceaccount/namespace"); ok {
		token.Namespace = fmt.Sprintf("%v", v)
	}

	// service account
	if v, ok := t.Get("kubernetes.io/serviceaccount/service-account.name"); ok {
		token.ServiceAccount = fmt.Sprintf("%v", v)
	}

	// secret
	if v, ok := t.Get("kubernetes.io/serviceaccount/secret.name"); ok {
		token.Secret = fmt.Sprintf("%v", v)
	}

	// return token
	return *&token, nil
}
