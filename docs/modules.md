# Modules

While writing test scripts, it is common to load part of different modules for usage throughout the script. In k6, it is possible to load four different kinds of modules.

- [Built-in modules](#built-in-modules)
- [Embedded modules](#embedded-modules)
- [Local modules](#local-modules)
- [Remote modules](#remote-modules)

Unlike in the real Python ecosystem, in the Starlark dialect (used in xk6-python), modules can be loaded with the [load](https://github.com/bazelbuild/starlark/blob/master/spec.md#load-statements) function instead of the import expression.

The first argument of the load function is the module name, followed by the symbols to be loaded (at least one). It is also possible to rename the symbol, for example to avoid collisions.

## Built-in modules

xk6-python provides many built-in modules for core functionalities.
For example, the `requests` client make HTTP requests against the system under test.

```python
load("check", "check")
load("requests", "get")

def default(_):
    resp = get("https://httpbin.test.k6.io/get")

    check(resp, {
        "is status 200": lambda r: r.status_code == 200,
    })
```

Built-in modules are implemented in go language and always available as part of the xk6-python extension. The name of the built-in modules does not contain a file extension.

## Embedded modules

xk6-python can provide embedded modules (currently only an example module). The embedded modules are always available, they are written in Python and embedded in the k6 executable binary. The name of the embedded modules does not contain a file extension.

```python
load("welcome", "hello")

def default(_):
    hello("Joe")
```

The embedded modules are used in exactly the same way as the built-in modules.

## Local modules

These modules are stored on the local filesystem, and accessed either through relative or absolute filesystem paths.

xk6-python doesn't support Python module resolution. File names for `load` must be fully specified (including file extension), such as `./helpers.py`.

```python
# my-test.py
load("./helpers.py", "someHelper")

def default(_):
  someHelper()
```

```python
# helpers.py
def someHelper():
  # ...
```

## Remote modules

These modules are accessed over `http(s)`, from a public source like GitHub, any CDN, or from any publicly accessible web server. The imported modules are downloaded and executed at runtime, making it extremely important to **make sure you trust the code before including it in a test script**.

For example, [k6pylib](https://github.com/grafana/k6pylib) is a set of k6 Python libraries available as remote `https` modules. They can be downloaded and imported as local modules or directly imported as remote modules.

```python
load("https://grafana.github.io/k6pylib/welcome.py", "hello")

def default(_):
    hello("Joe")
```

You can also build your custom Python libraries and distribute them via a public web hosting.

Remote modules can only be loaded using the `https` protocol. For development purposes, the `http` protocol is also allowed from the `127.0.0.1` host.

---

*This page is based on the original [Grafana k6 Modules](https://grafana.com/docs/k6/latest/using-k6/modules/) document.*
