"""
Example k6 test script to demonstrate the use of response data.
"""

load("requests", "get")

def default(_):
    """
    The default function defines the logic of the VUs.
    """
    resp = get("https://httpbin.test.k6.io/get")

    print(
        "url: '%s', reason: '%s', status_code: %d, ok: %s" %
        (resp.url, resp.reason, resp.status_code, resp.ok),
    )  # buildifier: disable=print

    print(resp.json())  # buildifier: disable=print
    print(resp.text)  # buildifier: disable=print
