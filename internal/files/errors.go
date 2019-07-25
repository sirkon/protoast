package files

import (
	"fmt"
)

var _ error = UnexpectedImportPath("")

type UnexpectedImportPath string

func (f UnexpectedImportPath) Error() string {
	return fmt.Sprintf("unexpected import %s", string(f))
}
