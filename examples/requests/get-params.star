"""
Example k6 test script to demonstrate the use of query parameters.
"""

load("requests", "get")

def default(_):
    """
    The default function defines the logic of the VUs.
    """
    resp = get("https://httpbin.test.k6.io/get", params = {"foo": "bar"})

    print(resp.json()["args"])  # buildifier: disable=print
