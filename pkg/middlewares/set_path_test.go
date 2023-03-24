package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatSetPathWorksOnRequests(t *testing.T) {
	mw := middlewares.MakeSetPathMiddleware("/xoxo")
	req := httptest.NewRequest(http.MethodGet, "/test?foo=bar", nil)
	mw(nil).ServeHTTP(nil, req)

	assert.Equal(t, "/xoxo?foo=bar", req.URL.RequestURI())
}
