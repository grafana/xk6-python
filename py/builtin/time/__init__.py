"""
The time module provides various time-related functions.
"""


def sleep(secs):
    """Suspend VU execution for the specified duration.

    :param secs: Duration, in seconds.

    Usage:

    .. code-block:: python

       load("time", "sleep")

       def default(_):
           sleep(0.5)
    """
