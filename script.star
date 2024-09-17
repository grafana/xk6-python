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
