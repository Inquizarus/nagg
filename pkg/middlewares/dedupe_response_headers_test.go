package middlewares_test

import (
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatDedupeResponseHeadersWorks(t *testing.T) {
	mw := middlewares.MakeDedupeResponseHeaders("x-foo", "x-bar")
	res := httptest.NewRecorder()

	res.Header().Add("x-foo", "foo1")
	res.Header().Add("x-foo", "foo2")

	res.Header().Add("x-bar", "bar1")
	res.Header().Add("x-bar", "bar2")

	mw(nil).ServeHTTP(res, nil)

	assert.Equal(t, "foo1", res.Header().Get("x-foo"))
	assert.Equal(t, "bar1", res.Header().Get("x-bar"))
}
