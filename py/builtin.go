package py

import (
	"go.k6.io/k6/js/modules"
	"go.starlark.net/starlark"
)

// BuiltinLoaderFunc loads a built-in (golang) module.
type BuiltinLoaderFunc func(string, *starlark.Thread, modules.VU) (starlark.StringDict, error)

var builtinRegistry = make(map[string]BuiltinLoaderFunc) //nolint:gochecknoglobals

// RegisterBuiltin registers a built-in module.
func RegisterBuiltin(loader BuiltinLoaderFunc, module string) {
	builtinRegistry[module] = loader
}

func loadBuiltin(module string, thread *starlark.Thread, vu modules.VU) (starlark.StringDict, bool, error) {
	loaderFunc, found := builtinRegistry[module]
	if !found {
		return nil, false, nil
	}

	globals, err := loaderFunc(module, thread, vu)
	if err != nil {
		return nil, true, err
	}

	return globals, true, nil
}
