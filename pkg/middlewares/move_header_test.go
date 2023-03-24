package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatMoveHeaderWorksForRequests(t *testing.T) {
	mw := middlewares.MakeMoveHeaderMiddleware("x-foo", "x-bar", "request")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-foo", "xoxo")
	mw(nil).ServeHTTP(nil, req)

	assert.Equal(t, "xoxo", req.Header.Get("x-bar"))
	assert.Equal(t, "", req.Header.Get("x-foo"))
}

func TestThatMoveHeaderWorksForResponses(t *testing.T) {
	mw := middlewares.MakeMoveHeaderMiddleware("x-foo", "x-bar", "response")
	res := httptest.NewRecorder()
	res.Header().Set("x-foo", "xoxo")
	mw(nil).ServeHTTP(res, nil)

	assert.Equal(t, "xoxo", res.Header().Get("x-bar"))
	assert.Equal(t, "", res.Header().Get("x-foo"))
}
