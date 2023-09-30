package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatSetResponseStatusCodeWorks(t *testing.T) {

	expectedStatusCode := http.StatusInternalServerError

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1", nil)
	res := httptest.NewRecorder()

	mw := middlewares.MakeSetResponseStatusCodeMiddleware(expectedStatusCode)

	mw(nil).ServeHTTP(res, req)

	assert.Equal(t, expectedStatusCode, res.Code)
}
