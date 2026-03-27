package core

import (
	"strings"
)

func (r *Registry) resolveName(scope, name string) (string, bool) {
	if strings.HasPrefix(name, ".") {
		_, ok := r.registry[name]
		return name, ok
	}
	if scope != "" && !strings.HasPrefix(scope, ".") {
		scope = "." + scope
	}
	for {
		cand := scope + "." + name
		if _, ok := r.registry[cand]; ok {
			return cand, true
		}
		if scope == "" {
			break
		}
		i := strings.LastIndex(scope, ".")
		if i <= 0 {
			scope = ""
		} else {
			scope = scope[:i]
		}
	}
	cand := "." + name
	_, ok := r.registry[cand]
	return cand, ok
}
