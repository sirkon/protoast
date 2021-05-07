package errors

import (
	"errors"
	"fmt"
)

// New вроппер для errors.New из стандартной библиотеки, вытащен сюда
// для того чтобы не распыляться по разным пакетам
func New(msg string) error {
	return errors.New(msg)
}

// Newf ошибка с форматированным текстом. Не является аналогом fmt.Errorf,
// т.к. не поддерживает флаг форматирования %w: данная функциональность
// возложена на функции Wrap(f)
func Newf(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}
