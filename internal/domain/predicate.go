package domain

import "net/http"

type predicates struct {
	path string
}

func (p *predicates) Path() string {
	return p.path
}

func (p *predicates) Matches(r *http.Request) bool {
	return p.Path() == r.URL.Path
}

func NewPredicates(path string) Predicates {
	return &predicates{path}
}
