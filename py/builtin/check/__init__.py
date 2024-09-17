"""
Checks validate boolean conditions in your test.

Testers often use checks to validate that the system is responding with the expected
content. For example, a check could validate that a POST request has a
response.status == 201, or that the body is of a certain size.

Checks are similar to what many testing frameworks call an assert, but failed checks
do not cause the test to abort or finish with a failed status. Instead, k6 keeps track
of the rate of failed checks as the test continues to run

Each check creates a rate metric. To make a check abort or fail a test, you can combine
it with a Threshold.
"""

from typing import Callable


def check(obj: object, verifications: list[Callable]) -> bool:
    """Run verifications on a given object.

    A verification is a test condition that can give a truthy or falsy result. Errors will
    be captured and will not interrupt the overall process.

    :param obj: the object to run the verifications on
    :param verifications: a dictionary, each key a description of the verification and each value
        a function to be called passing `obj`

    :return: True if all verifications returned True

    Usage:

    .. code-block:: python

       load("check", "check")
       load("requests", "get")

       def default(_):
           resp = get("https://httpbin.test.k6.io/get")

           check(resp, {
               "is status 200": lambda r: r.status_code == 200,
           })
    """
