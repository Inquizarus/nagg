package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatSetHeaderWorksForRequests(t *testing.T) {
	mw := middlewares.MakeSetHeaderMiddleware("x-test", "xoxo", "request")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(nil).ServeHTTP(nil, req)

	assert.Equal(t, "xoxo", req.Header.Get("x-test"))
}

func TestThatSetHeaderWorksForResponses(t *testing.T) {
	mw := middlewares.MakeSetHeaderMiddleware("x-test", "xoxo", "response")
	res := httptest.NewRecorder()
	mw(nil).ServeHTTP(res, nil)

	assert.Equal(t, "xoxo", res.Header().Get("x-test"))
}
