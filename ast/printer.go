package ast

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

// FileManager работа с записываемыми файлами
type FileManager interface {
	// Create возвращаем пару writer и функцию закрытия файла вместо одного io.WriteClose для удобства в тестах
	// Плюс, этот интерфейс будет использоваться только внутри приложения
	Create(name string) (writer io.Writer, closer func() error, err error)
}

// NewFileManager реализация управления файлами на уровне файловой системы
func NewFileManager(root string) FileManager {
	return rootFileManager(root)
}

var _ FileManager = rootFileManager("")

type rootFileManager string

// Create для реализации FileManager
func (r rootFileManager) Create(name string) (io.Writer, func() error, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, nil, err
	}
	return file, file.Close, nil
}

// NewPrinter конструктор принтера файлов
func NewPrinter(fm FileManager) *Printer {
	return &Printer{
		printed: map[string]struct{}{},
		plan:    map[string]*File{},
		fm:      fm,
	}
}

// Printer печать файлов
type Printer struct {
	printed map[string]struct{}
	plan    map[string]*File
	fm      FileManager
}

// Print печать файла и всех его зависимостей
func (p *Printer) Print(file *File) error {
	if _, ok := p.printed[file.Name]; ok {
		return nil
	}

	dest, closer, err := p.fm.Create(file.Name)
	if err != nil {
		return errors.Wrapf(err, "create f to print %s", file.Name)
	}
	defer func() {
		if err := closer(); err != nil {
			panic(errors.Wrapf(err, "close %s", file.Name))
		}
	}()

	if err := file.print(dest, p); err != nil {
		return errors.Wrapf(err, "print %s", file.Name)
	}
	p.printed[file.Name] = struct{}{}
	if _, ok := p.plan[file.Name]; ok {
		delete(p.plan, file.Name)
	}

	plan := make(map[string]*File, len(p.plan))
	for name, f := range p.plan {
		plan[name] = f
	}
	for _, f := range plan {
		if err := p.Print(f); err != nil {
			return errors.Wrapf(err, "print %s dependency", file.Name)
		}
	}

	return nil
}

// Plan добавить файл в план печати
func (p *Printer) Plan(f *File) {
	if _, ok := p.printed[f.Name]; ok {
		return
	}
	p.plan[f.Name] = f
}
