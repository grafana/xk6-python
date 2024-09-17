from typing_extensions import Self

"""
The metrics module provides functionality to create custom metrics of various types.
"""


class Metric:
    """A metric of one of the four custom metric types: counter, gauge, rate or trend.
    """

    def add(x: float | bool, tags: dict) -> Self:
        """Add a value to the metric.

        :param x: The value to add to the metric, either a number or a boolean. The boolean is interpreted as 0 (false) or 1 (true)
        :param tags: optional set of tags that will be tagged to the added data point (note that tags are added per data point and not for the entire metric).
        """


def counter(name: str) -> Metric:
    """Counter is an object for representing a custom cumulative counter metric.

    :param name: The name of the custom metric.

    :return: an initialized counter :py:class:`Metric`
    """


def gauge(name: str) -> Metric:
    """Gauge is an object for representing a custom metric holding only the latest value added.

    :param name: The name of the custom metric.

    :return: an initialized gauge :py:class:`Metric`
    """


def rate(name: str) -> Metric:
    """Rate is an object for representing a custom metric keeping track of the percentage of added values that are non-zero.

    :param name: The name of the custom metric.

    :return: an initialized rate :py:class:`Metric`
    """


def trend(name: str) -> Metric:
    """Trend is an object for representing a custom metric that allows for calculating different statistics on the added values (min, max, average or percentiles).

    :param name: The name of the custom metric.

    :return: an initialized trend :py:class:`Metric`
    """
