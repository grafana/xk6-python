"""
Example k6 test script.
"""

load("check", "check")

options = {}


def setup():
    data = {
        "foo": 123,
        "bar": "456",
    }
    return data


def teardown(data):
    pass


def _run(data, functions):
    result = check(data, functions)
    if result:
        print("==== yes! let's dance!")
    else:
        print("==== bad monkey :(")


def default(data):

    print("==== check example 1: should be ok")
    _run(
        data,
        {
        "foo is 123": lambda obj: obj["foo"] == 123,
        "string bar is positive": lambda obj: int(obj["bar"]) > 0,
        }
    )

    print("==== check example 2: should fail because comparison")
    _run(
        data,
        {
        "foo is 123": lambda obj: obj["foo"] == 999,  # bad value
        "string bar is positive": lambda obj: int(obj["bar"]) > 0,
        }
    )

    print("==== check example 3: should fail because crashing")
    _run(
        data,
        {
        "foo is 123": lambda obj: obj["foo"] == 123,
        "string bar is positive": lambda obj: obj["bar"] > 0,  # forgot the int()!!!
        }
    )
