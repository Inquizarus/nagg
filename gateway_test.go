package nagg_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/inquizarus/nagg"
	"github.com/inquizarus/nagg/pkg/logging"
	"github.com/inquizarus/rwapper/v2/pkg/servemuxwrapper"
	"github.com/stretchr/testify/assert"
)

func TestThatGatewayHandlerWorks(t *testing.T) {

	config, _ := nagg.JSONConfigFromFile("./testdata/gateway.json", nil)
	service := nagg.NewService(config)
	router := servemuxwrapper.New(nil)

	upstreamResponseBody := []byte(`{"foo":"bar"}`)

	router.HandlerFunc(http.MethodGet, "/api/example", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Add("x-upstream-header", "foobar")
		w.Header().Add("Content-Type", "application/json")
		w.Write(upstreamResponseBody)
		code, _ := strconv.Atoi(r.URL.Query().Get("status_code"))

		w.WriteHeader(code)
	})

	nagg.RegisterHTTPHandlers("/", router, service, logging.NewPlainLogger(nil, ""))

	server := httptest.NewServer(router)

	os.Setenv("HTTP_TEST_ADDRESS", server.URL+"/api/example")
	os.Setenv("HTTP_TEST_ADDRESS_NOT_FOUND", server.URL+"/api/not_found")

	tests := []struct {
		Destination               string
		ExpectedStatusCode        int
		ExpectedResponseBody      string
		ExpectedRequestHeaders    map[string]string
		ExpectedResponseHeaders   map[string]string
		ExpectedRequestParameters map[string]string
	}{
		{
			Destination:        server.URL + "/404",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Destination:        server.URL + "/404_no_upstream",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Destination:          server.URL + "/200?foo=bar&fizz=buzz",
			ExpectedStatusCode:   http.StatusOK,
			ExpectedResponseBody: string(upstreamResponseBody),
			ExpectedRequestHeaders: map[string]string{
				"x-global-header": "global",
			},
			ExpectedResponseHeaders: map[string]string{
				"x-post-header":     "post",
				"x-upstream-header": "foobar",
				"content-type":      "application/json",
			},
			ExpectedRequestParameters: map[string]string{
				"foo": "bar",
			},
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodGet, test.Destination, nil)

		req.URL.RawQuery = "status_code=" + fmt.Sprint(test.ExpectedStatusCode) + "&" + req.URL.RawQuery

		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assert.Equal(t, test.ExpectedStatusCode, res.Result().StatusCode, test.Destination)
		assert.Equal(t, test.ExpectedResponseBody, res.Body.String())

		for k, v := range test.ExpectedRequestHeaders {
			assert.Equal(t, v, req.Header.Get(k), test.Destination, k, v)
		}

		for k, v := range test.ExpectedResponseHeaders {
			assert.Equal(t, v, res.Header().Get(k), "incorrect value for response header "+k)
		}

		for k, v := range test.ExpectedRequestParameters {
			assert.True(t, req.URL.Query().Has(k))
			if v != "" {
				assert.Equal(t, v, req.URL.Query().Get(k), "incorrect value for request header "+k)
			}
		}
	}

	os.Unsetenv("HTTP_TEST_ADDRESS")
	os.Unsetenv("HTTP_TEST_ADDRESS_NOT_FOUND")
}
