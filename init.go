// Package xk6python contains xk6-python extension's entry point.
package xk6python

import "github.com/grafana/xk6-python/py"

func init() { //nolint:gochecknoinits
	py.Bootstrap()
}
