package metrics

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
	"go.starlark.net/starlark"

	"github.com/grafana/xk6-python/py/builtin/helpers"
)

type module struct {
	logger logrus.FieldLogger
	vu     modules.VU
}

// The Load function initializes the module and returns the module's global starlark symbols.
func Load(_ string, _ *starlark.Thread, vu modules.VU) (starlark.StringDict, error) {
	mod := &module{logger: vu.InitEnv().Logger, vu: vu}

	return helpers.ExportBuiltins(map[string]helpers.StarlarkFunction{
		"counter": mod.metric,
		"gauge":   mod.metric,
		"rate":    mod.metric,
		"trend":   mod.metric,
	}), nil
}

func (mod *module) metric(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var str string
	starlark.UnpackArgs(b.Name(), args, kwargs, "data", &str)
	if str == "" {
		return nil, fmt.Errorf("must provide a metric name, got: %s", args)
	}

	var m *NewMetric
	var err error
	switch b.Name() {
	case "counter":
		m, err = mod.newMetric(str, metrics.Counter)
	case "gauge":
		m, err = mod.newMetric(str, metrics.Gauge)
	case "rate":
		m, err = mod.newMetric(str, metrics.Rate)
	case "trend":
		m, err = mod.newMetric(str, metrics.Trend)
	default:
		err = errors.New("unknown metric type")
	}

	return m, err
}

func (mod *module) newMetric(name string, t metrics.MetricType) (*NewMetric, error) {
	initEnv := mod.vu.InitEnv()
	if initEnv == nil {
		return nil, errors.New("metrics must be declared in the init context")
	}
	valueType := metrics.Default
	m, err := initEnv.Registry.NewMetric(name, t, valueType)
	if err != nil {
		return nil, err
	}
	metric := &NewMetric{
		metric: m,
		vu:     mod.vu,
		logger: mod.logger,
	}

	return metric, nil
}
