// The greeting package is an example builtin module.
package greeting

import (
	"github.com/grafana/xk6-python/py/builtin/helpers"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
	"go.starlark.net/starlark"
)

type module struct {
	logger logrus.FieldLogger
}

// The Load function initializes the module and returns the module's global starlark symbols.
func Load(_ string, _ *starlark.Thread, vu modules.VU) (starlark.StringDict, error) {
	mod := &module{logger: vu.InitEnv().Logger}

	return helpers.ExportBuiltins(map[string]helpers.StarlarkFunction{
		"hi": mod.hi,
	}), nil
}

func (mod *module) hi(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	name := "World"

	if args.Len() > 0 {
		str, ok := starlark.AsString(args.Index(0))
		if ok {
			name = str
		}
	}

	mod.logger.Infof("Hi, %s!", name)

	return starlark.None, nil
}
