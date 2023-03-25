package middlewares

import (
	"net/http"

	"github.com/MicahParks/keyfunc"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/inquizarus/nagg/pkg/httptools"
	"github.com/inquizarus/nagg/pkg/logging"
)

func MakeCheckJWTValidityByJWKSURL(url string, log logging.Logger, client *http.Client) func(http.Handler) http.Handler {

	if client == nil {
		client = http.DefaultClient
	}

	if log == nil {
		log = logging.DefaultLogger
	}

	jwks, err := keyfunc.Get(url, keyfunc.Options{
		Client: client,
	})

	if err != nil {
		log.Errorf("Failed to get the JWKS from the given URL, %s", err)
	}

	return makeCheckJWTWithJWKSMiddleware(jwks, log)

}

func makeCheckJWTWithJWKSMiddleware(jwks *keyfunc.JWKS, log logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if jwks == nil {
				log.Error("could not check JWT as JWKS is nil")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b64Token, err := httptools.ExtractBearerTokenFromRequest(r)

			if err != nil {
				log.Errorf("could not retrieve bearer token from request, %s", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(b64Token, jwks.Keyfunc)

			if err != nil {
				log.Errorf("could not parse jwt token, %s", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				log.Debugf("jwt token in request was not valid, %s", token.Claims.Valid())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if nil != next {
				next.ServeHTTP(w, r)
			}
		})
	}
}
