"""
Example k6 test script to demonstrate the use of form post data.
"""

load("requests", "post")

def default(_):
    """
    The default function defines the logic of the VUs.
    """
    resp = post(
        "https://httpbin.test.k6.io/post",
        data = {"foo": "bar"},
    )

    print(resp.json()["form"])  # buildifier: disable=print
