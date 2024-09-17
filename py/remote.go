package py

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.k6.io/k6/js/modules"
	"go.starlark.net/starlark"
)

func loadRemote(module string, thread *starlark.Thread, vu modules.VU) (starlark.StringDict, bool, error) {
	if !strings.HasPrefix(module, "https:") && !strings.HasPrefix(module, "http://127.0.0.1") {
		return nil, false, nil
	}

	req, err := http.NewRequestWithContext(vu.Context(), http.MethodGet, module, nil)
	if err != nil {
		return nil, true, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, true, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, true, fmt.Errorf("%w: %s", errRemote, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, true, err
	}

	globals, err := starlark.ExecFileOptions(fileOptions, thread, module, data, predeclared(module, module))
	if err != nil {
		return nil, true, err
	}

	return globals, true, nil
}

var errRemote = errors.New("remote error")
