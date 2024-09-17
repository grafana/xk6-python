package py

import (
	"errors"
	"io/fs"
	"os"

	"go.k6.io/k6/js/modules"
	"go.starlark.net/starlark"
)

var filesystemRegistry = make([]fs.FS, 0, 1) //nolint:gochecknoglobals

// RegisterFilesystem registers all embedded modules from a filesystem.
func RegisterFilesystem(filesystem fs.FS) {
	filesystemRegistry = append(filesystemRegistry, filesystem)
}

// RegisterSubFilesystem registers all embedded modules from a direcory of a filesystem.
func RegisterSubFilesystem(filesystem fs.FS, dir string) error {
	sub, err := fs.Sub(filesystem, dir)
	if err != nil {
		return err
	}

	RegisterFilesystem(sub)

	return nil
}

func loadFilesystem(module string, thread *starlark.Thread, vu modules.VU) (starlark.StringDict, bool, error) {
	for _, filesystem := range filesystemRegistry {
		globals, found, err := loadFS(filesystem, module, thread, vu)
		if found {
			return globals, true, err
		}
	}

	return nil, false, nil
}

func loadFS(fsys fs.FS, module string, thread *starlark.Thread, _ modules.VU) (starlark.StringDict, bool, error) {
	file := module

	for _, ext := range []string{".py", ".star"} {
		if _, err := fs.Stat(fsys, module+ext); err == nil {
			file += ext
			break
		}
	}

	if _, err := fs.Stat(fsys, file); errors.Is(err, os.ErrNotExist) { //nolint:forbidigo
		return nil, false, nil
	}

	data, err := fs.ReadFile(fsys, file)
	if err != nil {
		return nil, true, err
	}

	globals, err := starlark.ExecFileOptions(fileOptions, thread, module, data, predeclared(file, module))
	if err != nil {
		return nil, true, err
	}

	return globals, true, nil
}
