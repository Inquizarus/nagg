package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatRedirectWorksOnRequests(t *testing.T) {
	mw := middlewares.MakeRedirectMiddleware(http.StatusTemporaryRedirect, "http://localhost")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	res := httptest.NewRecorder()
	mw(nil).ServeHTTP(res, req)
	assert.Equal(t, http.StatusTemporaryRedirect, res.Code)
	assert.Equal(t, "http://localhost", res.Header().Get("location"))
}
