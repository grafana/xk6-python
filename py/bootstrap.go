package py

import (
	"embed"

	"go.k6.io/k6/js/modules"

	"github.com/grafana/xk6-python/py/builtin/check"
	"github.com/grafana/xk6-python/py/builtin/greeting"
	"github.com/grafana/xk6-python/py/builtin/group"
	"github.com/grafana/xk6-python/py/builtin/json"
	"github.com/grafana/xk6-python/py/builtin/metrics"
	"github.com/grafana/xk6-python/py/builtin/requests"
	"github.com/grafana/xk6-python/py/builtin/time"
)

//go:embed embedded
var embedded embed.FS

// Bootstrap initializes the extension. Register the extension as a k6 JavaScript module.
// In the case of the run command, if the argument has .py or .star file extension,
// it rewrites the argument to stdin and sends a bootstrap JavaScript code to stdin.
func Bootstrap() {
	redirectStdin()

	// register k6 extension
	modules.Register("k6/x/python", newExtension())

	// register built-in modules
	RegisterBuiltin(check.Load, "check")
	RegisterBuiltin(greeting.Load, "greeting")
	RegisterBuiltin(group.Load, "group")
	RegisterBuiltin(json.Load, "json")
	RegisterBuiltin(metrics.Load, "metrics")
	RegisterBuiltin(requests.Load, "requests")
	RegisterBuiltin(time.Load, "time")

	// register embedded source modules
	_ = RegisterSubFilesystem(embedded, "embedded")
}
