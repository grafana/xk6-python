// Implement 'check' functionality

package check

import (
	"fmt"
	"time"

	"github.com/grafana/xk6-python/py/builtin/helpers"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/js/modules/k6"
	"go.k6.io/k6/metrics"
	"go.starlark.net/starlark"
)

type module struct {
	vu modules.VU
	logger logrus.FieldLogger
}

// The Load function initializes the module and returns the module's global starlark symbols.
func Load(_ string, _ *starlark.Thread, vu modules.VU) (starlark.StringDict, error) {
	mod := &module{logger: vu.InitEnv().Logger, vu: vu}

	return helpers.ExportBuiltins(map[string]helpers.StarlarkFunction{
		"check": mod.check,
	}), nil
}

func (mod *module) check(th *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	state := mod.vu.State()
	if state == nil {
		return starlark.False, k6.ErrCheckInInitContext
	}

	// validate the arguments and get useful objects
	mod.logger.Debugf("Running check, args len=%d", args.Len())
	if args.Len() != 2 {
		return nil, fmt.Errorf("Bad arguments len=%d (expected 2)", args.Len())
	}
	obj, verifications := args[0], args[1]

	// check verifications is a dict
	verificationsDict, ok := verifications.(*starlark.Dict)
	if !ok {
		return nil, fmt.Errorf("Error: 'verifications' is not a dictionary")
	}

	// prepare the metric tags
	commonTagsAndMeta := state.Tags.GetCurrentValues()

	ctx := mod.vu.Context()
	allItemsOk := true
	currentTime := time.Now()

	// iterate verifications key and function, validate types, and call the function
	for idx, item := range verificationsDict.Items() {
		thisItemOk := true

		key, value := item[0], item[1]
		keyString, ok := key.(starlark.String)
		if !ok {
			return nil, fmt.Errorf("Error: 'verifications' key of item %d is not a string", idx)
		}
		valueFunction, ok := value.(starlark.Callable)
		if !ok {
			return nil, fmt.Errorf("Error: 'verifications' value of item %d is not a function", idx)
		}

		mod.logger.Debugf("Running function of item %d: %s", idx, keyString)
		result, err := starlark.Call(th, valueFunction, starlark.Tuple{obj}, nil)
		if err == nil {
			mod.logger.Debugf("Function returned: %s", result)
			resultBool := starlark.Bool(result.Truth())
			if !resultBool  {
				thisItemOk = false
			}
		} else {
			mod.logger.Errorf("Function crashed: %s", err)
			thisItemOk = false
		}

		if !thisItemOk {
			allItemsOk = false
		}

		tags := commonTagsAndMeta.Tags
		if state.Options.SystemTags.Has(metrics.TagCheck) {
			tags = tags.With("check", string(keyString))
		}

		sample := metrics.Sample{
			TimeSeries: metrics.TimeSeries{
				Metric: state.BuiltinMetrics.Checks,
				Tags: tags,
			},
			Time: currentTime,
		}
		if thisItemOk {
			sample.Value = 1
		}

		metrics.PushIfNotDone(ctx, state.Samples, sample)
	}

	if allItemsOk {
		// all functions call succeded ok and returned True
		mod.logger.Debug("All functions checked ok")
	}
	return starlark.Bool(allItemsOk), nil
}
