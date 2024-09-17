package py

import "go.starlark.net/starlark"

func predeclared(file string, module string) starlark.StringDict {
	if len(module) == 0 {
		module = "__main__"
	}

	if len(file) == 0 {
		file = "<stdin>"
	}

	return starlark.StringDict{
		"__name__": starlark.String(module),
		"__file__": starlark.String(file),
	}
}

// samplePredeclared is used only to determine is a name is predeclared or not.
var samplePredeclared = predeclared("", "") //nolint:gochecknoglobals

func isPredeclared(name string) bool {
	_, found := samplePredeclared[name]

	return found
}
