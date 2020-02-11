// Code generated by ldetool --go-string component.lde. DO NOT EDIT.

package component

import (
	"strings"
)

// Component ...
type Component struct {
	Rest   string
	Dir    string
	Object string
}

// Extract ...
func (p *Component) Extract(line string) (bool, error) {
	p.Rest = line
	var pos int

	// Take until ':' as Dir(string)
	pos = strings.IndexByte(p.Rest, ':')
	if pos >= 0 {
		p.Dir = p.Rest[:pos]
		p.Rest = p.Rest[pos+1:]
	} else {
		return false, nil
	}

	// Take the rest as Object(string)
	p.Object = p.Rest
	p.Rest = p.Rest[len(p.Rest):]
	return true, nil
}