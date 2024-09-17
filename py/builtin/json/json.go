package json

import (
	"go.k6.io/k6/js/modules"
	"go.starlark.net/lib/json"
	"go.starlark.net/starlark"
)

// The Load function initializes the module and returns the module's global starlark symbols.
// This is a very thin wrapper to load starlark included json methods,
func Load(_ string, _ *starlark.Thread, vu modules.VU) (starlark.StringDict, error) {
	members := json.Module.Members
	// Alias the methods to the names commonly used in python
	// See https://docs.python.org/3/library/json.html
	members["loads"] = json.Module.Members["decode"]
	members["dumps"] = json.Module.Members["encode"]
	return members, nil
}
