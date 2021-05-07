package errors

import "fmt"

// Wrap прямой аналог fmt.Errorf("%s: %w", msg, err)
func Wrap(err error, msg string) error {
	err.Error() // чтобы вылетать при аннотации пустой ошибки
	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf аналогично Wrap, но с форматированной аннотацией ошибки
func Wrapf(err error, format string, a ...interface{}) error {
	err.Error()
	return fmt.Errorf(format+": %w", append(a, err)...)
}
