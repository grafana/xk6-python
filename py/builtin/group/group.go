// Implement 'group' functionality

package group

import (
	"fmt"
	"time"

	"github.com/grafana/xk6-python/py/builtin/helpers"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/js/modules/k6"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/metrics"
	"go.starlark.net/starlark"
)

type module struct {
	vu     modules.VU
	logger logrus.FieldLogger
}

// The Load function initializes the module and returns the module's global starlark symbols.
func Load(_ string, _ *starlark.Thread, vu modules.VU) (starlark.StringDict, error) {
	mod := &module{logger: vu.InitEnv().Logger, vu: vu}

	return helpers.ExportBuiltins(map[string]helpers.StarlarkFunction{
		"group": mod.group,
	}), nil
}

func (mod *module) group(th *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	state := mod.vu.State()
	if state == nil {
		return starlark.False, k6.ErrCheckInInitContext
	}

	// validate the arguments and get useful objects
	mod.logger.Debugf("Running group, args len=%d", args.Len())
	if args.Len() != 2 {
		return nil, fmt.Errorf("Bad arguments len=%d (expected 2)", args.Len())
	}
	_name, _function := args[0], args[1]

	// validate the type of the arguments
	name, ok := _name.(starlark.String)
	if !ok {
		return nil, fmt.Errorf("Error: 'name' argument must be a string")
	}
	function, ok := _function.(starlark.Callable)
	if !ok {
		return nil, fmt.Errorf("Error: 'function' argument must be a callable")
	}

	oldGroupName, _ := state.Tags.GetCurrentValues().Tags.Get(metrics.TagGroup.String())
	// TODO: what are we doing if group is not tagged
	newGroupName, err := lib.NewGroupPath(oldGroupName, name.GoString())
	if err != nil {
		return nil, err
	}

	shouldUpdateTag := state.Options.SystemTags.Has(metrics.TagGroup)
	if shouldUpdateTag {
		state.Tags.Modify(func(tagsAndMeta *metrics.TagsAndMeta) {
			tagsAndMeta.SetSystemTagOrMeta(metrics.TagGroup, newGroupName)
		})
	}
	defer func() {
		if shouldUpdateTag {
			state.Tags.Modify(func(tagsAndMeta *metrics.TagsAndMeta) {
				tagsAndMeta.SetSystemTagOrMeta(metrics.TagGroup, oldGroupName)
			})
		}
	}()

	startTime := time.Now()
	result, err := starlark.Call(th, function, nil, nil)
	endTime := time.Now()

	ctx := mod.vu.Context()
	ctm := state.Tags.GetCurrentValues()
	metrics.PushIfNotDone(ctx, state.Samples, metrics.Sample{
		TimeSeries: metrics.TimeSeries{
			Metric: state.BuiltinMetrics.GroupDuration,
			Tags:   ctm.Tags,
		},
		Time:     endTime,
		Value:    metrics.D(endTime.Sub(startTime)),
		Metadata: ctm.Metadata,
	})

	return result, err
}
