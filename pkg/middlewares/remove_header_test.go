package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatRemoveHeaderRequests(t *testing.T) {
	mw := middlewares.MakeRemoveHeaderMiddleware("x-test", "request")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(nil).ServeHTTP(nil, req)

	assert.Equal(t, "", req.Header.Get("x-test"))
}

func TestThatRemoveHeaderResponses(t *testing.T) {
	mw := middlewares.MakeRemoveHeaderMiddleware("x-test", "response")
	res := httptest.NewRecorder()
	mw(nil).ServeHTTP(res, nil)

	assert.Equal(t, "", res.Header().Get("x-test"))
}
