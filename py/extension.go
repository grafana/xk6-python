// Package py contains xk6-python module implementation.
package py

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"go.k6.io/k6/js/modules"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

var _ modules.Module = (*extension)(nil)

type extension struct {
	prog *starlark.Program

	compileOnce sync.Once
}

func newExtension() *extension {
	ext := new(extension)

	return ext
}

//nolint:forbidigo
func (ext *extension) compile(filename string) error {
	script, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return err
	}

	_, prog, err := starlark.SourceProgramOptions(fileOptions, filename, script, isPredeclared)
	if err != nil {
		return err
	}

	ext.prog = prog

	return nil
}

func (ext *extension) init(filename string, throw func(error)) {
	ext.compileOnce.Do(func() {
		if err := ext.compile(filename); err != nil {
			throw(err)
		}
	})
}

//nolint:forbidigo
func (ext *extension) NewModuleInstance(vu modules.VU) modules.Instance { //nolint:varnamelen
	mod := &instance{ //nolint:exhaustruct
		vu: vu,
	}

	ext.init(os.Getenv(envScript), mod.throw)

	mod.init(ext.prog)

	return mod
}

const (
	envScript = "XK6_PYTHON_SCRIPT"

	nameSetup    = "setup"
	nameDefault  = "default"
	nameTeardown = "teardown"
	nameOptions  = "options"
)

var errMissingSymbol = errors.New("missing symbol")

var fileOptions = &syntax.FileOptions{ //nolint:gochecknoglobals
	Set:             true,
	While:           true,
	TopLevelControl: true,
	GlobalReassign:  true,
}
