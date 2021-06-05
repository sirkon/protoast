package ast

import "errors"

var _ error = ErrorTypeNotFound("")

// ErrorTypeNotFound ошибка возвращаемая если тип не был найден
type ErrorTypeNotFound string

func (e ErrorTypeNotFound) Error() string {
	return string(e) + ": type not found"
}

// IsErrorTypeNotFound проверка, что данная ошибка является обёрткой для ErrorTypeNotFound
func IsErrorTypeNotFound(err error) bool {
	var target ErrorTypeNotFound
	// не получается использовать errors.Is из-за его ориентации на производительность
	return errors.As(err, &target)
}

var _ error = ErrorPackageMissingFiles{}

// ErrorPackageMissingFiles ошибка указывающая на отсутствие proto-файлов в пакете
type ErrorPackageMissingFiles struct{}

func (ErrorPackageMissingFiles) Error() string {
	return "missing proto files to form a package"
}
