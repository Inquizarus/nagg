package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatSetRequestParameterWorks(t *testing.T) {
	mw := middlewares.MakeSetRequestParameterMiddleware("foo", "bar")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	mw(nil).ServeHTTP(nil, req)

	assert.Equal(t, "bar", req.URL.Query().Get("foo"))
}
