package liner

import (
	"io"

	"github.com/sirkon/go-format"
)

func New(w io.Writer) Liner {
	return Liner{
		w: w,
	}
}

type Liner struct {
	w io.Writer
}

func (p Liner) Line(line string, a ...interface{}) {
	if _, err := io.WriteString(p.w, format.Formatp(line, a...)); err != nil {
		panic(err)
	}
	if _, err := io.WriteString(p.w, "\n"); err != nil {
		panic(err)
	}
}

func (p Liner) Newl() {
	if _, err := io.WriteString(p.w, "\n"); err != nil {
		panic(err)
	}
}
