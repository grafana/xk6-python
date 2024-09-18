# xk6-python

**Write k6 tests in Python**

> [!IMPORTANT]
> `xk6-python` is not an official k6 extension but was a **Grafana Hackathon #10** project!
> Active development is not planned in the near future.
> The purpose of making it public is to assess demand (number of stars on the repo).

The xk6-python k6 extension allows writing k6 tests in ([a dialect of](https://github.com/google/starlark-go/blob/master/doc/spec.md)) the Python programming language. The tests are executed by a built-in Python interpreter ([starlark-go](https://github.com/google/starlark-go)), so there is no need for Python installation.

```python file=script.star
"""
Example k6 test script.
"""

load("check", "check")
load("requests", "get")
load("time", "sleep")

options = {
    "vus": 5,
    "duration": "5s",
    "thresholds": {
        "checks": ["rate>=0.99"],
    },
}

def default(_):
    resp = get("https://httpbin.test.k6.io/get")

    check(resp, {
        "is status 200": lambda r: r.status_code == 200,
    })

    sleep(0.5)
```

```bash
./k6 run script.py
```

_This is not a toy, but a Proof of Concept for full-fledged k6 Python language support._

Check out the [documentation](https://grafana.github.io/xk6-python/).

## Motivation

Python is quite a popular programming language these days. According to the [TIOBE Programming Community Index](https://www.tiobe.com/tiobe-index/) 2024, Python has secured the top position, beating C++, C, Java, and JavaScript.

Even though [k6 intentionally only supports one programming language](https://k6.io/blog/why-k6-does-not-introduce-multiple-scripting-languages/), it is worth considering making an exception for Python.

## Why Starlark?

The [starlark-go](https://github.com/google/starlark-go) package is a pure go implementation of the [Starlark python dialect](https://github.com/google/starlark-go/blob/master/doc/spec.md). Its use does not require external dependencies, such as the installation of CPython. It doesn't even require the use of [cgo](https://go.dev/wiki/cgo). This enables the portability of the k6 executable binary and simplifies its use in the cloud.

## Is it Python?

Yes and no.

Yes, because starlark is a python dialect. The **language and syntax is python**.

No, because the **usual python ecosystem cannot be used**:
- no python module system
- packages that can be installed with pip cannot be used.

This is an embedded python interpreter that does not make the python ecosystem available in the same way as the k6 JavaScript interpreter does not make the Node.js reason system available. So only the python language and syntax can be used plus the built-in modules and the local and remote python/starlark modules.

## Features

The xk6-python development is in the PoC phase, but it can be used to write real k6 tests. Currently implemented main features:

- supports the k6 options object as a Python dict
- fully supports k6 lifecycle functions (setup, default, teardown)
- supports the use of built-in, remote and local modules
- partial support of the most important k6 APIs

## Download

You can download pre-built k6 binaries from [Releases](https://github.com/grafana/xk6-python/releases/) page.

## Status

**xk6-python** is currently in _Proof of Concept_ status, but can already be used to run real k6 tests written in Python.

Check the [documentation](https://grafana.github.io/xk6-python/) for available modules.

## Contribute

If you would like to contribute, start by reading [Contributing Guidelines](https://grafana.github.io/xk6-python/CONTRIBUTING.html).

## k6 version

xk6-python is currently compatible with k6 v0.52.0. It is possible that it works with other versions, but it has not been tested.
