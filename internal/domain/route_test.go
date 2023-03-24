package domain_test

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/inquizarus/nagg/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestThatRouteMatchWorks(t *testing.T) {
	tests := []struct {
		route    domain.Route
		request  *http.Request
		expected bool
	}{
		{
			route: domain.NewRoute("example", "", domain.NewPredicates("/example"), nil),
			request: &http.Request{
				URL: &url.URL{
					Path: "/example",
				},
			},
			expected: true,
		},
		{
			route: domain.NewRoute("example", "", domain.NewPredicates("/example"), nil),
			request: &http.Request{
				URL: &url.URL{
					Path: "/",
				},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.route.Matches(test.request))
	}
}

func TestThatRouteAddressWorks(t *testing.T) {
	os.Setenv("TestThatRouteMatchWorks_Address", "http://localhost:1234")
	tests := []struct {
		route    domain.Route
		request  *http.Request
		expected string
	}{
		{
			route: domain.NewRoute("example", "env:TestThatRouteMatchWorks_Address", domain.NewPredicates("/example"), nil),
			request: &http.Request{
				URL: &url.URL{
					Path: "/example",
				},
			},
			expected: "http://localhost:1234",
		},
		{
			route: domain.NewRoute("example", "http://localhost:5678", domain.NewPredicates("/example"), nil),
			request: &http.Request{
				URL: &url.URL{
					Path: "/",
				},
			},
			expected: "http://localhost:5678",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.route.Address())
	}
	os.Unsetenv("TestThatRouteMatchWorks_Address")
}
