package files

import (
	"fmt"
)

var _ error = UnexpectedImportPath("")

// UnexpectedImportPath представление ошибки об отсутствующих данных по импорту
type UnexpectedImportPath string

func (f UnexpectedImportPath) Error() string {
	return fmt.Sprintf("unexpected import %s", string(f))
}
