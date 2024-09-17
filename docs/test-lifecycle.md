# Test lifecycle

In the lifecycle of a k6 test, a script always runs through these stages in the same order:

1. Code in the `init` context prepares the script, loading files, load modules, and defining the test _lifecycle functions_. **Required**.
2. The `setup` function runs, setting up the test environment and generating data. _Optional._
3. VU code runs in the `default` or scenario function, running for as long and as many times as the `options` define. **Required**.
4. The `teardown` function runs, postprocessing data and closing the test environment. _Optional._

```python
# 1. init code

def setup():
  # 2. setup code
  return {} # data arg for default() and teardown()

def default(data):
  # 3. VU code
  pass

def teardown(data):
  # 4. teardown code
  pass
```

**Lifecycle functions**

Except for init code, each stage occurs in a _lifecycle function_, a function called in a specific sequence in the k6 runtime.

## Overview of the lifecycle stages

For examples and implementation details of each stage, refer to the subsequent sections.

| Test stage      | Purpose                                                       | Example                                                                                 | Called                                                        |
| --------------- | ------------------------------------------------------------- | --------------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| **1. init**     | Read local files, load modules, declare lifecycle functions | Open JSON file, Load module                                                           | Once per VU                                                 |
| **2. Setup**    | Set up data for processing, share data among VUs              | Call API to start test environment                                                      | Once                                                          |
| **3. VU code**  | Run the test function, usually `default`                      | Make https requests, validate responses                                                 | Once per iteration, as many times as the test options require |
| **4. Teardown** | Process result of setup code, stop test environment           | Validate that setup had a certain result, send webhook notifying that test has finished | Once \*                                                     |


\* If the `Setup` function ends abnormally (e.g throws an error), the `teardown()` function isn't called. Consider adding logic to the `setup()` function to handle errors and ensure proper cleanup.

## The init stage

**The init stage is required**.
Before the test runs, k6 needs to initialize the test conditions.
To prepare the test, code in the `init` context runs once per VU.

Some operations that might happen in `init` include the following:

- Load modules
- Load files from the local file system
- Configure the test for all `options`
- Define lifecycle functions for the VU, `setup`, and `teardown` stages

**All code that is outside of a lifecycle function is code in the `init` context**.
Code in the `init` context _always executes first_.


```python
# init context: loading modules
load("metrics", "trend")
load("requests", http_get = "get")

# init context: define k6 options
options = {
  "vus": 10,
  "duration": "30s",
}

# init context: global variables
customTrend = trend("oneCustomMetric")

# init context: define custom function
def myCustomFunction():
```

Separating the `init` stage from the VU stage removes irrelevant computation from VU code, which improves k6 performance and makes test results more reliable.
One limitation of `init` code is that it **cannot** make HTTP requests.
This limitation ensures that the `init` stage is reproducible across tests (the response from protocol requests is dynamic and unpredictable)

## The VU stage

Scripts must contain, at least, a _scenario function_ that defines the logic of the VUs.
The code inside this function is _VU code_.
Typically, VU code is inside the `default` function, but it can also be inside the function defined by a scenario (see subsequent section for an example).

```python
def default(_):
  # do things here...
```

**VU code runs over and over through the test duration.**
VU code can make HTTP requests, emit metrics, and generally do everything you'd expect a load test to do.
The only exceptions are the jobs that happen in the `init` context.

- VU code _does not_ load files from your local filesystem.
- VU code _does not_ import any other modules.

Again, instead of VU code, init code does these jobs.

### The default function life-cycle

A VU executes the `default()` function from start to end in sequence.
Once the VU reaches the end of the function, it loops back to the start and executes the code all over.

As part of this "restart" process, k6 resets the VU.
Cookies are cleared, and TCP connections might be torn down (depending on your test configuration options).

## Setup and teardown stages

Like `default`, `setup` and `teardown` functions must be functions.
But unlike the `default` function, k6 calls `setup` and `teardown` only once per test.

- `setup` is called at the beginning of the test, after the init stage but before the VU stage.
- `teardown` is called at the end of a test, after the VU stage (`default` function).

You can call the full k6 API in the setup and teardown stages, unlike the init stage.
For example, you can make HTTP requests:

```python
load("requests", "get")

def setup():
  res = get("https://httpbin.test.k6.io/get")
  return res.json()

def teardown(data):
  print(data)

def default(data):
  print(data)
```

### Skip setup and teardown execution

You can skip the execution of setup and teardown stages using the options `--no-setup` and
`--no-teardown`.


```bash
k6 run --no-setup --no-teardown ...
```

### Use data from setup in default and teardown

Again, let's have a look at the basic structure of a k6 test:

```python
# 1. init code

def setup():
  # 2. setup code

def default(data):
  # 3. VU code

def teardown(data):
  # 4. teardown code
```

You might have noticed the function signatures of the `default()` and `teardown()` functions take an argument, referred to here as `data`.

Here's an example of passing some data from the setup code to the VU and teardown stages:

```python
def setup():
  return { "v" : 1 }

def default(data):
  print(data)

def teardown(data)
  if data["v"] != 1:
    print("incorrect data: ",data)
```

For example, with the data returned by the `setup()` function, you can:

- Give each VU access to an identical copy of the data
- Postprocess the data in `teardown` code

However, there are some restrictions.

- You can pass only data (i.e. JSON) between `setup` and the other stages.
  You cannot pass functions.
- If the data returned by the `setup()` function is large, it will consume more memory.
- You cannot manipulate data in the `default()` function, then pass it to the `teardown()` function.

It's best to think that each stage and each VU has access to a fresh "copy" of whatever data the `setup()` function returns.

![Diagram showing data getting returned by setup, then used (separately) by default and teardown functions](https://grafana.com/media/docs/k6-oss/lifecycle.png)

It would be extremely complicated and computationally intensive to pass mutable data between all VUs and then to teardown, especially in distributed setups.
This would go against a core k6 goal: the same script should be executable in multiple modes.

---

*This page was created with minimal modifications to the original [Grafana k6 Test lifecycle](https://grafana.com/docs/k6/latest/using-k6/test-lifecycle/) document, converting the JavaScript examples to Python.*
