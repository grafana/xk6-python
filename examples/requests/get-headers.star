"""
Example k6 test script to demonstrate the use of request headers.
"""

load("requests", "get")

def default(_):
    """
    The default function defines the logic of the VUs.
    """
    resp = get("https://httpbin.test.k6.io/get", headers = {"foo": "bar"})

    print(resp.json()["headers"])  # buildifier: disable=print
