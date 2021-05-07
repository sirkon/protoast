package protoast

import (
	"fmt"
	"text/scanner"

	"github.com/sirkon/protoast/internal/errors"
)

var _ error = errorPosition{}

type errorPosition struct {
	pos scanner.Position
	err error
}

func (e errorPosition) Error() string {
	return fmt.Sprintf("%s %s", e.pos, e.err)
}

func errPos(pos scanner.Position, err error) error {
	return errorPosition{
		pos: pos,
		err: err,
	}
}

func errPosf(pos scanner.Position, format string, a ...interface{}) error {
	return errPos(pos, errors.Newf(format, a...))
}
