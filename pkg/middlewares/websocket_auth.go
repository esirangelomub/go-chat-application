package middlewares

import (
	"github.com/go-chi/jwtauth"
	"net/http"
)

func Verifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verify(ja, jwtauth.TokenFromQuery, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)(next)
	}
}
