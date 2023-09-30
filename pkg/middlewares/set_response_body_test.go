package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatSetResponseBodyWorks(t *testing.T) {

	expectedBody := "{\"foo\":\"bar\"}"
	reader := strings.NewReader(expectedBody)

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1", nil)
	res := httptest.NewRecorder()

	mw := middlewares.MakeSetResponseBodyMiddleware(reader)

	mw(nil).ServeHTTP(res, req)

	assert.Equal(t, expectedBody, res.Body.String())
}

func TestThatSetResponseBodyWorksMultipleTimes(t *testing.T) {

	expectedBody := "{\"foo\":\"bar\"}"
	reader := strings.NewReader(expectedBody)

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1", nil)
	res1 := httptest.NewRecorder()
	res2 := httptest.NewRecorder()

	mw := middlewares.MakeSetResponseBodyMiddleware(reader)

	mw(nil).ServeHTTP(res1, req)
	mw(nil).ServeHTTP(res2, req)

	assert.Equal(t, expectedBody, res1.Body.String())
	assert.Equal(t, expectedBody, res2.Body.String())
}
