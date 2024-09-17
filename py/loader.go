package py

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.k6.io/k6/js/modules"
	"go.starlark.net/starlark"
)

type loader struct {
	vu    modules.VU
	cache map[string]*cacheEntry
}

type cacheEntry struct {
	globals starlark.StringDict
	err     error
}

func newLoader(vu modules.VU) *loader {
	return &loader{
		vu:    vu,
		cache: make(map[string]*cacheEntry),
	}
}

func (l *loader) load(invoker *starlark.Thread, module string) (starlark.StringDict, error) {
	entry, found := l.cache[module]
	if entry != nil {
		return entry.globals, entry.err
	}

	if found {
		// request for package whose "load in progress".
		return nil, fmt.Errorf("cycle in load graph")
	}

	// indicate "load in progress".
	l.cache[module] = nil

	thread := &starlark.Thread{Name: "exec " + module, Load: invoker.Load}

	entry = new(cacheEntry)

	entry.globals, found, entry.err = loadRemote(module, thread, l.vu)
	if !found {
		if ext := filepath.Ext(module); len(ext) == 0 {
			entry.globals, found, entry.err = loadBuiltin(module, thread, l.vu)
			if !found {
				entry.globals, found, entry.err = loadFilesystem(module, thread, l.vu)
			}
		} else {
			entry.globals, found, entry.err = loadFS(os.DirFS("."), module, thread, l.vu) //nolint:forbidigo
		}
	}

	if !found {
		entry.err = fmt.Errorf("%w: %s", errModuleNotFound, module)
	}

	l.cache[module] = entry

	return entry.globals, entry.err
}

var errModuleNotFound = errors.New("module not found")
