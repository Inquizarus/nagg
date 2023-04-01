package domain_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/inquizarus/nagg/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestThatPredicateMatchWorks(t *testing.T) {
	tests := []struct {
		predicate domain.Predicates
		request   *http.Request
		expected  bool
	}{
		{
			predicate: domain.NewPredicates("/example"),
			request: &http.Request{
				URL: &url.URL{
					Path: "/example",
				},
			},
			expected: true,
		},
		{
			predicate: domain.NewPredicates("/example"),
			request: &http.Request{
				URL: &url.URL{
					Path: "/",
				},
			},
			expected: false,
		},
		{
			predicate: domain.NewPredicates("/example"),
			request: &http.Request{
				URL: &url.URL{
					Path: "/example/ostkaka",
				},
			},
			expected: false,
		},
		{
			predicate: domain.NewPredicates("/example*"),
			request: &http.Request{
				URL: &url.URL{
					Path: "/example/ostkaka",
				},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.predicate.Matches(test.request))
	}
}
