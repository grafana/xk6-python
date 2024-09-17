package helpers

import "go.starlark.net/starlark"

type StarlarkFunction func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error)

func ExportBuiltins(funcs map[string]StarlarkFunction) starlark.StringDict {
	globals := make(starlark.StringDict)

	for name, fn := range funcs {
		globals[name] = starlark.NewBuiltin(name, fn)
	}

	return globals
}
