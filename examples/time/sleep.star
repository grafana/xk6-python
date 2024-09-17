"""
Example k6 test script to demonstrate the use of the sleep function.
"""

load("time", "sleep")

def default(_):
    """
    The default function defines the logic of the VUs.
    """
    print("Before sleep")  # buildifier: disable=print
    sleep(1.5)
    print("After sleep")  # buildifier: disable=print
