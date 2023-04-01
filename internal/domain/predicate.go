package domain

import (
	"net/http"
	"strings"
)

type predicates struct {
	path string
}

func (p *predicates) Path() string {
	return p.path
}

func (p *predicates) Matches(r *http.Request) bool {

	if p.pathIsPrefix() {
		return strings.HasPrefix(r.URL.Path, p.path[:len(p.path)-1])
	}

	return p.Path() == r.URL.Path
}

func (p *predicates) pathIsPrefix() bool {
	return strings.HasSuffix(p.path, "*")
}

func NewPredicates(path string) Predicates {
	return &predicates{path}
}
