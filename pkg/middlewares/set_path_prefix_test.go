package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatSetPathPrefixWorksOnRequests(t *testing.T) {
	mw := middlewares.MakeSetRequestPathPrefixMiddleware("/xoxo")
	req := httptest.NewRequest(http.MethodGet, "/test?foo=bar", nil)
	mw(nil).ServeHTTP(nil, req)

	assert.Equal(t, "/xoxo/test?foo=bar", req.URL.RequestURI())
}
