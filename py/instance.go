package py

import (
	"encoding/json"
	"fmt"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.starlark.net/starlark"
)

type instance struct {
	vu      modules.VU
	loader  *loader
	thread  *starlark.Thread
	globals starlark.StringDict

	exports modules.Exports
}

var _ modules.Instance = (*instance)(nil)

func (inst *instance) init(prog *starlark.Program) {
	inst.thread = new(starlark.Thread)

	logger := inst.vu.InitEnv().Logger.WithField("source", "console")

	inst.thread.Print = func(_ *starlark.Thread, msg string) {
		logger.Info(msg)
	}

	inst.loader = newLoader(inst.vu)

	inst.thread.Load = inst.loader.load

	globals, err := prog.Init(inst.thread, predeclared(prog.Filename(), ""))
	if err != nil {
		inst.throw(err)
	}

	inst.globals = globals

	if err := inst.initExports(); err != nil {
		inst.throw(err)
	}
}

func (inst *instance) throw(err error) {
	common.Throw(inst.vu.Runtime(), err)
}

func (inst *instance) callSetup() (interface{}, error) {
	inst.vu.State().Logger.Debug("Calling Setup")

	value, err := inst.global(nameSetup)
	if err != nil {
		return nil, err
	}

	svalue, err := starlark.Call(inst.thread, value, nil, nil)
	if err != nil {
		return nil, err
	}

	return svalue.String(), nil
}

func (inst *instance) callTeardown(data interface{}) error {
	inst.vu.State().Logger.Debug("Calling Teardown")

	fn, err := inst.global(nameTeardown)
	if err != nil {
		return err
	}

	return inst.callWithData(fn, data)
}

func (inst *instance) callDefault(data interface{}) error {
	inst.vu.State().Logger.Debug("Calling Default")

	fn, err := inst.global(nameDefault)
	if err != nil {
		return err
	}

	return inst.callWithData(fn, data)
}

func (inst *instance) callWithData(fn starlark.Value, data interface{}) error {
	var args starlark.Tuple

	if data != nil {
		sdata, err := starlark.EvalOptions(fileOptions, inst.thread, fn.String()+"(data)", data, nil)
		if err != nil {
			return err
		}

		args = starlark.Tuple{sdata}
	} else {
		args = starlark.Tuple{starlark.None}
	}

	_, err := starlark.Call(inst.thread, fn, args, nil)

	return err
}

func (inst *instance) initExports() error {
	toValue := inst.vu.Runtime().ToValue

	inst.exports = modules.Exports{
		Default: nil,
		Named:   make(map[string]interface{}),
	}

	if inst.globals.Has(nameSetup) {
		inst.exports.Named[nameSetup] = toValue(inst.callSetup)
	}

	if inst.globals.Has(nameDefault) {
		inst.exports.Default = toValue(inst.callDefault)
	}

	if inst.globals.Has(nameTeardown) {
		inst.exports.Named[nameTeardown] = toValue(inst.callTeardown)
	}

	if inst.globals.Has(nameOptions) {
		value, err := inst.global(nameOptions)
		if err != nil {
			return err
		}

		var opts map[string]interface{}

		if err := json.Unmarshal([]byte(value.String()), &opts); err != nil {
			return err
		}

		inst.exports.Named[nameOptions] = toValue(opts)
	}

	return nil
}

func (inst *instance) Exports() modules.Exports {
	return inst.exports
}

func (inst *instance) global(name string) (starlark.Value, error) {
	value, ok := inst.globals[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", errMissingSymbol, name)
	}

	return value, nil
}
