package metrics

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
	"go.starlark.net/starlark"
)

type NewMetric struct {
	metric *metrics.Metric
	vu     modules.VU
	logger logrus.FieldLogger
}

func metricAdd(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x starlark.Value
	var tags *starlark.Dict
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &x, &tags); err != nil {
		return nil, err
	}
	br := b.Receiver()
	m, ok := br.(*NewMetric)
	if !ok {
		return nil, fmt.Errorf("function must be called on a metric, got: %s", reflect.TypeOf(br))
	}

	value := float64(0)
	if x.Type() == "bool" {
		if x.Truth() {
			value = 1
		} else {
			value = 0
		}
	} else if v, err := starlark.AsInt32(x); err == nil {
		value = float64(v)
	} else {
		return nil, fmt.Errorf("cannot add to metric '%s': %w", m.metric.Name, err)
	}

	ctm := m.vu.State().Tags.GetCurrentValues()
	if tags != nil {
		for _, tag := range tags.Items() {
			key, ok := starlark.AsString(tag[0])
			if !ok {
				return nil, fmt.Errorf("tag key '%s' is not a string", tag[0])
			}
			val, ok := starlark.AsString(tag[1])
			if !ok {
				return nil, fmt.Errorf("tag value '%s' is not a string", tag[1])
			}
			ctm.SetTag(key, val)
		}
	}

	sample := metrics.Sample{
		TimeSeries: metrics.TimeSeries{
			Metric: m.metric,
			Tags:   ctm.Tags,
		},
		Time:     time.Now(),
		Metadata: ctm.Metadata,
		Value:    value,
	}

	metrics.PushIfNotDone(m.vu.Context(), m.vu.State().Samples, sample)

	return m, nil
}

func (m *NewMetric) Attr(name string) (starlark.Value, error) {
	if name == "name" {
		return starlark.String(m.metric.Name), nil
	}
	if name == "add" {
		return starlark.NewBuiltin(
			"add",
			metricAdd,
		).BindReceiver(m), nil
	}
	return nil, fmt.Errorf("unknown function on metric, got: %s", name)
}

func (m *NewMetric) AttrNames() []string {
	return []string{"add", "name"}
}

// Implement the [starlark.Value] interface
func (m *NewMetric) String() string { return m.metric.Name }

func (*NewMetric) Type() string { return "metric" }

func (*NewMetric) Freeze() {}

func (*NewMetric) Truth() starlark.Bool { return true }

func (*NewMetric) Hash() (uint32, error) { return 0, errors.New("unhashable") }
