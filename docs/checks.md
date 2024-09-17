# Checks

Checks validate boolean conditions in your test. Testers often use checks to validate that the system is responding with the expected content. For example, a check could validate that a POST request has a `response.status_code == 201`, or that the body is of a certain size.

Checks are similar to what many testing frameworks call an _assert_, but **failed checks do not cause the test to abort or finish with a failed status**. Instead, k6 keeps track of the rate of failed checks as the test continues to run.

Each check creates a `checks` rate metric. To make a check abort or fail a test, you can combine it with a `threshold`.

## Check for HTTP response code

Checks are great for codifying assertions relating to HTTP requests and responses.
For example, this snippet makes sure the HTTP response code is a 200:

```python
load("check", "check")
load("requests", "get")

def default():
    resp = get("https://httpbin.test.k6.io/get")

    check(resp, {
        "is status 200": lambda r: r.status_code == 200,
    })

```

## See percentage of checks that passed

When a script includes checks, the summary report shows how many of the tests' checks passed:

```shell
$ k6 run script.js

  ...
    ✓ is status 200

  ...
  checks.........................: 100.00% ✓ 1        ✗ 0
  data_received..................: 11 kB   12 kB/s
```

In this example, note that the check "is status 200" succeeded 100% of the times it was called.

## Fail a load test using checks

When a check fails, the script will continue executing successfully and will not return a 'failed' exit status. Checks are nice for codifying assertions, but unlike `thresholds`, `checks` do not affect the exit status of k6.

If you need the whole test to fail based on the results of a check, you have to combine checks with thresholds. This is particularly useful in specific contexts, such as integrating k6 into your CI pipelines or receiving alerts when scheduling your performance tests.

```javascript
load("check", "check")
load("requests", "get")
load("time", "sleep")

options = {
    "vus": 5,
    "duration": "5s",
    "thresholds": {
        "checks": ["rate>0.9"],
    },
}

def default():
    resp = get("https://httpbin.test.k6.io/get")

    check(resp, {
        "is status 200": lambda r: r.status_code == 200,
    })

    sleep(0.5)
```

In this example, the threshold is configured on the `checks` metric, establishing that the rate of successful checks is higher than `90%`.
