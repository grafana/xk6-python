# Extend functionality

The functionality of xk6-python can be extended with remote modules written in Starlark and built-in modules written in go.

## By using remote modules

The functionality of xk6-python can be extended most easily with remote modules written in Starlark language. It is not necessary to modify the source code of the xk6-python extension, nor is it necessary to build a new k6 executable.

A Starlark module should be created and deployed on a static web server accessible via `https` protocol. The remote modules can be used on the developer's machine with the `http` protocol by specifying the host `127.0.0.1`.

```python
load("https://grafana.github.io/k6pylib/welcome.py", "hello")

def default(_):
    hello("Joe")
```

See more in the [Remote Modules](https://grafana.github.io/xk6-python/modules.html#remote-modules) section.

## By using built-in modules

The built-in modules should be written in the go programming language. A custom k6 build is required to integrate them.

The built-in module is loaded by a loader function implemented by the module.

```go
type BuiltinLoaderFunc func(string, *starlark.Thread, modules.VU) (starlark.StringDict, error)
```

The module loader function should be registered using the `RegisterBuiltin()` function during startup.

```go
func RegisterBuiltin(loader BuiltinLoaderFunc, module string)
```

There are two options for placing the source code of the module:

1. By modifying the source of xk6-pyhon. The module should be placed in a subdirectory of the [py/builtin](https://github.com/grafana/xk6-python/tree/main/py/builtin) directory and registered in the `Bootstrap()` function.

2. By creating a k6 extension. This k6 extension should be integrated into k6 and the module registration should be done in the `init()` function.

For example, it is worth looking at the [built-in modules of xk6-python](https://github.com/grafana/xk6-python/tree/main/py/builtin).
