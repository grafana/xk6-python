"""
Example k6 test script to demonstrate the use of JSON post data.
"""

load("requests", "post")

def default(_):
    """
    The default function defines the logic of the VUs.
    """
    resp = post(
        "https://httpbin.test.k6.io/post",
        json = {"foo": "bar"},
    )

    print(resp.json()["json"])  # buildifier: disable=print
