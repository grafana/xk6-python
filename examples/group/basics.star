"""
Example k6 test script.
"""

load("check", "check")
load("group", "group")
load("requests", "get")

options = {}


def default(_):

    group("a simple get to google", lambda: get("http://google.com"))

    def _work():
        resp = get("http://google.com")
        check(resp, {
            "is ok": lambda r: r.status_code == 200,
            "not empty": lambda r: len(r.text) > 0,
        })

    group("a more comprehensive get to google", _work)

    group("silly math", lambda: 1 * 2 * 3)
