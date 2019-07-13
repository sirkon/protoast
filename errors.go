package prototypes

import (
	"fmt"
	"text/scanner"

	"github.com/pkg/errors"
)

var _ error = errorPosition{}

type errorPosition struct {
	pos scanner.Position
	err error
}

func (e errorPosition) Error() string {
	return fmt.Sprintf("%s %s", e.pos, e.err)
}

// errPos возвращает позиционную ошибку
func errPos(pos scanner.Position, err error) error {
	return errorPosition{
		pos: pos,
		err: err,
	}
}

// errPosf позиционная ошибка с форматом
func errPosf(pos scanner.Position, format string, a ...interface{}) error {
	return errPos(pos, errors.Errorf(format, a...))
}
