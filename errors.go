package prototypes

import (
	"fmt"
	"text/scanner"

	"github.com/pkg/errors"
)

var _ error = errPos{}

type errPos struct {
	pos scanner.Position
	err error
}

func (e errPos) Error() string {
	return fmt.Sprintf("%s %s", e.pos, e.err)
}

// ErrPos возвращает позиционную ошибку
func ErrPos(pos scanner.Position, err error) error {
	return errPos{
		pos: pos,
		err: err,
	}
}

// ErrPosf позиционная ошибка с форматом
func ErrPosf(pos scanner.Position, format string, a ...interface{}) error {
	return ErrPos(pos, errors.Errorf(format, a...))
}
