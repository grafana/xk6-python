"""
Example k6 test script to demonstrate the use of a custom metric.
"""

load("metrics", "gauge")

myGauge = gauge("my_gauge")

def default(_):
    """
    The default function defines the logic of the VUs.
    """
    myGauge.add(1)
    myGauge.add(43,  { "tag1": 'value', "tag2": 'value2' })
