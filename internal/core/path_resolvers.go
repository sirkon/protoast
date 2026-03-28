package core

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/sirkon/protoast/v2/internal/errors"
)

// PathResolver returns absolute path for the given path in protobuf import statement.
type PathResolver interface {
	String() string
	Resolve(path string) (string, error)
}

type pathResolverProtoc struct {
	root string
}

func newPathResolverProtoc() (*pathResolverProtoc, error) {
	var out bytes.Buffer
	command := exec.Command("which", "protoc")
	command.Stdout = &out
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return nil, errors.Wrap(err, "look for protoc executable path")
	}

	dir, _ := filepath.Split(out.String())
	dir, _ = filepath.Split(strings.TrimSuffix(dir, string(filepath.Separator)))

	include := filepath.Join(dir, "include")

	// test for google/protobuf/any.proto
	testPath := filepath.Join(include, "google", "protobuf", "any.proto")
	if _, err := os.Stat(testPath); err != nil {
		return nil, errors.Wrapf(err, "check for google/protobuf/any.proto near the protoc in %q", include)
	}

	return &pathResolverProtoc{
		root: include,
	}, nil
}

func (p *pathResolverProtoc) String() string {
	return "protoc distribution"
}

func (p *pathResolverProtoc) Resolve(path string) (string, error) {
	fullPath := filepath.Join(p.root, path)

	if _, err := os.Stat(fullPath); err != nil {
		return "", errors.Wrap(err, "check computed path")
	}

	return fullPath, nil
}

type pathResolverRoot struct {
	root string
}

func newPathResolverRoot(root string) (*pathResolverRoot, error) {
	stat, err := os.Stat(root)
	if err != nil {
		return nil, errors.Wrap(err, "check root path")
	}

	if !stat.IsDir() {
		return nil, errors.New("root path is not a directory")
	}

	return &pathResolverRoot{root: root}, nil
}

func (r *pathResolverRoot) String() string {
	return fmt.Sprintf("schema at %q", r.root)
}

func (r *pathResolverRoot) Resolve(path string) (string, error) {
	fullPath := filepath.Join(r.root, path)
	if _, err := os.Stat(fullPath); err != nil {
		return "", errors.Wrap(err, "check computed path")
	}

	return fullPath, nil
}

type PathResolversBuilder struct {
	isProtoc bool
	roots    []string
}

func Resolvers() *PathResolversBuilder {
	return &PathResolversBuilder{}
}

func (b *PathResolversBuilder) WithProtoc() *PathResolversBuilder {
	b.isProtoc = true
	return b
}

func (b *PathResolversBuilder) WithRoot(roots ...string) *PathResolversBuilder {
	b.roots = append(b.roots, roots...)
	return b
}

func (b *PathResolversBuilder) Build() ([]PathResolver, error) {
	roots := slices.Clone(b.roots)
	sort.Strings(roots)
	slices.Compact(roots)

	var result []PathResolver

	if b.isProtoc {
		res, err := newPathResolverProtoc()
		if err != nil {
			return nil, errors.Wrap(err, "setup protoc distribution resolver")
		}

		result = append(result, res)
	}

	for _, root := range roots {
		resolver, err := newPathResolverRoot(root)
		if err != nil {
			return nil, errors.Wrap(err, "setup schema resolver over "+root)
		}

		result = append(result, resolver)
	}

	return result, nil
}
