"""
Example k6 test script to demonstrate the use of various custom metrics.
"""

load("metrics", "counter", "gauge", "rate", "trend")

myCounter = counter("my_counter") # A metric that cumulatively sums added values.
myGauge = gauge("my_gauge") # A metric that stores the min, max and last values added to it.
myRate = rate("my_rate") # A metric that tracks the percentage of added values that are non-zero.
myTrend = trend("my_trend") # A metric that calculates statistics on the added values (min, max, average, and percentiles).

def default(_):
    """
    The default function defines the logic of the VUs.
    """
    myCounter.add(1)
    myGauge.add(9)
    myTrend.add(231)
    myRate.add(17)

    # Set of tags that will be tagged to the added data point (note that tags are added per data point and not for the entire metric).
    tags = { "tag1": 'value', "tag2": 'value2' }

    myTrend.add(2, tags)
    myCounter.add(45, tags)
    myCounter.add(2, {"other_tag": "another_value"})
