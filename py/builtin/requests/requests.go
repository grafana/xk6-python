// The requests package contains HTTP requests module.
package requests

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/grafana/xk6-python/py/builtin/helpers"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"

	"go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type module struct {
	vu     modules.VU
	client *http.Client
	logger logrus.FieldLogger
}

// The Load function initializes the module and returns the module's global starlark symbols.
func Load(_ string, _ *starlark.Thread, vu modules.VU) (starlark.StringDict, error) {
	mod := &module{
		vu:     vu,
		logger: vu.InitEnv().Logger,
		client: &http.Client{Transport: newTransport(vu)},
	}

	return helpers.ExportBuiltins(map[string]helpers.StarlarkFunction{
		"request": mod.request,
		"head":    mod.head,
		"get":     mod.get,
		"post":    mod.post,
		"put":     mod.put,
		"patch":   mod.patch,
		"delete":  mod.delete,
	}), nil
}

func (mod *module) execute(thread *starlark.Thread, method string, url string, kwargs []starlark.Tuple) (starlark.Value, error) {
	req, err := mod.newRequest(thread, method, url, kwargs)
	if err != nil {
		return nil, err
	}

	log := mod.logger.WithFields(logrus.Fields{"url": req.URL, "method": req.Method})

	log.Debug("HTTP request")

	resp, err := mod.client.Do(req)
	if err != nil {
		return nil, err
	}

	return mod.newResponse(resp)
}

func (mod *module) request(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	mod.logger.Debug(builtin.Name())

	if err := helpers.RequireArgsLen(builtin.Name(), args, 2); err != nil {
		return nil, err
	}

	method, err := helpers.RequireStringArg(builtin.Name(), args, 0, "method")
	if err != nil {
		return nil, err
	}

	url, err := helpers.RequireStringArg(builtin.Name(), args, 1, "url")
	if err != nil {
		return nil, err
	}

	return mod.execute(thread, method, url, kwargs)
}

func (mod *module) head(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	url, err := helpers.RequireStringArg(builtin.Name(), args, 0, "url")
	if err != nil {
		return nil, err
	}

	return mod.execute(thread, http.MethodHead, url, kwargs)
}

func (mod *module) get(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	url, err := helpers.RequireStringArg(builtin.Name(), args, 0, "url")
	if err != nil {
		return nil, err
	}

	return mod.execute(thread, http.MethodGet, url, kwargs)
}

func (mod *module) post(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	url, err := helpers.RequireStringArg(builtin.Name(), args, 0, "url")
	if err != nil {
		return nil, err
	}

	return mod.execute(thread, http.MethodPost, url, kwargs)
}

func (mod *module) put(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	url, err := helpers.RequireStringArg(builtin.Name(), args, 0, "url")
	if err != nil {
		return nil, err
	}

	return mod.execute(thread, http.MethodPut, url, kwargs)
}

func (mod *module) patch(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	url, err := helpers.RequireStringArg(builtin.Name(), args, 0, "url")
	if err != nil {
		return nil, err
	}

	return mod.execute(thread, http.MethodPatch, url, kwargs)
}

func (mod *module) delete(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	url, err := helpers.RequireStringArg(builtin.Name(), args, 0, "url")
	if err != nil {
		return nil, err
	}

	return mod.execute(thread, http.MethodDelete, url, kwargs)
}

func (mod *module) newRequest(thread *starlark.Thread, method string, loc string, kwargs []starlark.Tuple) (*http.Request, error) {
	opts := helpers.OptionalArgs(kwargs)

	body, contentType, err := newRequestBody(thread, opts)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(mod.vu.Context(), method, loc, body)
	if err != nil {
		return nil, err
	}

	if headers, has := opts["headers"]; has {
		dict, ok := headers.(*starlark.Dict)
		if !ok {
			return nil, fmt.Errorf("%w: headers arg must be a dict", errInvalidType)
		}

		for _, item := range dict.Items() {
			req.Header.Add(toString(item.Index(0)), toString(item.Index(1)))
		}
	}

	if len(req.Header.Get("Content-Type")) == 0 && len(contentType) != 0 {
		req.Header.Set("Content-Type", contentType)
	}

	req.URL.RawQuery = newQuery(opts["params"])

	return req, nil
}

func (mod *module) newResponse(resp *http.Response) (starlark.Value, error) {
	dict := make(starlark.StringDict)

	dict["url"] = starlark.String(resp.Request.URL.String())
	dict["reason"] = starlark.String(resp.Status)
	dict["status_code"] = starlark.MakeInt(resp.StatusCode)
	dict["headers"] = mod.newHeaders(resp.Header)
	dict["ok"] = starlark.Bool(resp.StatusCode < http.StatusBadRequest)

	defer resp.Body.Close()

	// TODO handle Content-Encoding response header
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dict["text"] = starlark.String(body)

	fn := func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		return starlark.Call(thread, json.Module.Members["decode"], starlark.Tuple{starlark.String(body)}, nil)
	}

	dict["json"] = starlark.NewBuiltin("json", fn)

	// TODO implement other response fields

	return starlarkstruct.FromStringDict(starlarkstruct.Default, dict), nil
}

func (mod *module) newHeaders(header http.Header) *starlark.Dict {
	dict := starlark.NewDict(len(header))

	for key := range header {
		dict.SetKey(starlark.String(key), starlark.String(header.Get(key)))
		// TODO implement custom case insensitive dict?
		dict.SetKey(starlark.String(strings.ToLower(key)), starlark.String(header.Get(key)))
	}

	return dict
}

func newQuery(params starlark.Value) string {
	if params == nil || params == starlark.None {
		return ""
	}

	switch value := params.(type) {
	case *starlark.Dict:
		query := make(url.Values, value.Len())

		for _, tuple := range value.Items() {
			key, _ := starlark.AsString(tuple.Index(0))

			val, ok := starlark.AsString(tuple.Index(1))
			if !ok {
				val = tuple.Index(1).String()
			}

			query.Add(key, val)
		}

		return query.Encode()
	}

	return ""
}

func newRequestBody(thread *starlark.Thread, opts starlark.StringDict) (io.Reader, string, error) {
	if data, has := opts["json"]; has {
		encoded, err := starlark.Call(thread, json.Module.Members["encode"], starlark.Tuple{data}, nil)
		if err != nil {
			return nil, "", err
		}

		str, _ := starlark.AsString(encoded)

		return strings.NewReader(str), "application/json; charset=utf-8", nil
	}

	if data, has := opts["data"]; has {
		var str string

		switch value := data.(type) {
		case *starlark.Dict:
			params := url.Values{}
			for _, item := range value.Items() {
				params.Set(toString(item.Index(0)), toString(item.Index(1)))
			}

			str = params.Encode()

		case *starlark.List:
			params := url.Values{}
			var item starlark.Value

			for iter := value.Iterate(); iter.Next(&item); {
				tuple, ok := item.(starlark.Tuple)
				if !ok {
					continue
				}

				params.Set(toString(tuple.Index(0)), toString(tuple.Index(1)))
			}

			str = params.Encode()

		case starlark.Bytes:
			str = string(value)
		case starlark.String:
			str = string(value)
		default:
			return nil, "", fmt.Errorf("%w: %T", errInvalidType, data)
		}

		return strings.NewReader(str), "application/x-www-form-urlencoded", nil
	}

	return nil, "", nil
}

func toString(value starlark.Value) string {
	str, ok := starlark.AsString(value)
	if !ok {
		str = value.String()
	}

	return str
}

var errInvalidType = errors.New("invalid type")
