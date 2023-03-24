package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatRemoveRequestParameterWorks(t *testing.T) {
	mw := middlewares.MakeRemoveRequestParameterMiddleware("foo")
	req := httptest.NewRequest(http.MethodGet, "/test?foo=bar", nil)
	mw(nil).ServeHTTP(nil, req)

	assert.Equal(t, "", req.URL.Query().Get("foo"))
}
