// The time package is the implementation of the builtin time module.
package time

import (
	"time"

	"github.com/grafana/xk6-python/py/builtin/helpers"
	"go.k6.io/k6/js/modules"
	"go.starlark.net/starlark"
)

type module struct {
	vu modules.VU
}

// The Load function initializes the module and returns the module's global starlark symbols.
func Load(_ string, _ *starlark.Thread, vu modules.VU) (starlark.StringDict, error) {
	mod := &module{vu: vu}

	return helpers.ExportBuiltins(map[string]helpers.StarlarkFunction{
		"sleep": mod.sleep,
	}), nil
}

func (mod *module) sleep(_ *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var secsValue starlark.Value

	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 1, &secsValue); err != nil {
		return nil, err
	}

	var secs float64

	if v, ok := starlark.AsFloat(secsValue); ok {
		secs = v
	} else {
		var v int

		if err := starlark.AsInt(secsValue, &v); err != nil {
			return nil, err
		}

		secs = float64(v)
	}

	ctx := mod.vu.Context()
	timer := time.NewTimer(time.Duration(secs * float64(time.Second)))

	select {
	case <-timer.C:
	case <-ctx.Done():
		timer.Stop()
	}

	return starlark.None, nil
}
