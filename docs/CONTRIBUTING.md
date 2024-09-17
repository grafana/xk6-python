# Contributing Guidelines

Thank you for your interest in contributing to xk6-python!

Before you begin, make sure to familiarize yourself with the [Code of Conduct](CODE_OF_CONDUCT.md). If you've previously contributed to other open source project, you may recognize it as the classic [Contributor Covenant](https://contributor-covenant.org/).

## Filing issues

Don't be afraid to file issues! Nobody can fix a bug we don't know exists, or add a feature we didn't think of.

1. **Ensure the bug was not already reported** by searching on GitHub under [Issues](https://github.com/grafana/xk6-python/issues).

2. If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/grafana/xk6-python/issues/new). Be sure to include a **title and clear description**, as much relevant information as possible.


The worst that can happen is that someone closes it and points you in the right direction.

## Contributing code

If you'd like to contribute code to xk6-python, you will need [Go](https://go.dev/doc/install), [golangci-lint](https://github.com/golangci/golangci-lint) and [xk6](https://github.com/grafana/xk6).

The starlark REPL tool can be useful for experimentation:

```
go install go.starlark.net/cmd/starlark@latest
```

This is the basic procedure.

1. Find an issue you'd like to fix. If there is none already, or you'd like to add a feature, please open one, and we can talk about how to do it. Out of respect for your time, please start a discussion regarding any bigger contributions in a GitHub Issue **before** you get started on the implementation.

2. Create a fork and open a feature branch - `feature/my-cool-feature` is the classic way to name these, but it really doesn't matter.

3. Create a pull request!

4. We will discuss implementation details until everyone is happy, then a maintainer will merge it.

## Implement a module

The new module can be implemented in two ways:

 - built-in module using the go programming language
 - embedded module using the starlark programming language
 - local module using the starlark programming language

Regardless of the way of implementation, the modules are used in the same way. The builtin module takes precedence if there is a builtin and an embedded module with the same name.

### Builtin module

The builtin module must be written in the go programming language. Each module is placed in a separate go package within the [py/builtin](https://github.com/grafana/xk6-python/tree/main/py/builtin) directory. The package must have a public loader function to load the module's global starlark symbols.

```go
type BuiltinLoaderFunc func(moduleName string, thread *starlark.Thread, vu modules.VU) (starlark.StringDict, error)
```

The builtin module must be registered in the `Bootstrap()` function (in the [py/bootstrap.go](https://github.com/grafana/xk6-python/tree/main/py/bootstrap.go) file) using a `registerBuiltin()` function call.

### Embedded module

The embedded modules are to be implemented in starlark (Python). Each module is placed in a separate file within the [py/embedded](https://github.com/grafana/xk6-python/tree/main/py/embedded) directory. The file extension can be .py or .star, the advantage of .star is that using the [bazel-stack-vscode](https://marketplace.visualstudio.com/items?itemName=StackBuild.bazel-stack-vscode) Visual Studio Code extension, quite good editor support is available for the starlark language.

Embedded modules are loaded without using a file extension. This way, the embedded module can later be rewritten as a builtin module.

No registration is required for modules placed in the [py/embedded](https://github.com/grafana/xk6-python/tree/main/py/embedded) directory. This folder will be embedded in the go code of the k6 extension and the modules here will be automatically available in the test script.

If justified, additional embedded filesystems can be added to the search path of embeddedd modules. Additional embedded filesystems can be registered in the `Bootstrap()` function (in the [py/bootstrap.go](https://github.com/grafana/xk6-python/tree/main/py/bootstrap.go) file) using a `registerFilesystem()` or `registerSubFilesystem` function call.

### Remote module

The remote modules are to be implemented in starlark (Python). The remote modules must be deployed on a static web server and can be accessed from there using the https protocol.

The [grafana/k6pylib](https://github.com/grafana/k6pylib) GitHub repository is an example remote module collection. The repository is automatically deployed using GitHub Pages.

The easiest way to try out remote modules is to use the [grafana/k6pylib](https://github.com/grafana/k6pylib) repository. You simply need to create a Python/starlark module and merge it into the `main` branch.

### Local module

The local module is only included in this list for the sake of completeness. The local modules are not part of the xk6-python extension, they are located locally next to the test script, in files with the extension `.py` or `.star`. To load the local module, you must always specify the file extension.

## Typical tasks

This section describes the typical tasks of contributing code.

If the [cdo](https://github.com/szkiba/cdo) tool is installed, the tasks can be easily executed. Otherwise, you'll need to type or copy the commands here.

### lint - Run the linter

The `golangci-lint` tool is used for static analysis of the source code.
It is advisable to run it before committing the changes.

```bash
golangci-lint run
```

[lint]: <#lint---run-the-linter>

### test - Run the tests

Run all tests and collect coverage data.

```bash
go test -count 1 -race -coverprofile=coverage.txt ./...
```

[test]: <#test---run-the-tests>

### coverage - View the test coverage report

Requires
: [test]

```bash
go tool cover -html=coverage.txt
```

### build - Build k6 executable binary

Extensions can be integrated into the k6 executable using the xk6 CLI tool.

```bash
xk6 build --with github.com/grafana/xk6-python=.
```

[build]: <#build---build-k6-binary>

### run - Run Python test script with k6

With the xk6 CLI tool, python test scripts can be run from the repository base directory without building k6.


```bash
xk6 run ${1:-script.py}
```

### readme - Update exmaple in README.md

Using the [mdcode] tool, the example script in README.md can be updated from the script.py file.

```bash
mdcode update
```

[mdcode]: <https://github.com/szkiba/mdcode>

### docs-install - Install documentation tools

Create a Python virtual environment, activate, and install the needed dependencies (for the first time only).

```bash
python3 -m venv env
source env/bin/activate
pip install -r docs/requirements.txt
```

[docs-install]: <#docs-install---install-documentation-tools>

### docs - Generate documentation

Generate project documentation in `build/docs` directory.

Requires
: [docs-install]

```bash
source env/bin/activate
sphinx-build -b html docs build/docs
```

### ci - Run all ci-relevant tasks

Run all the tasks that will run in the Continuous Integration system.

Requires
: [lint], [test], [build]
