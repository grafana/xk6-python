"""
Run code inside a group. Groups are used to organize results in a test.
"""

from typing_extensions import Callable, Any


def group(name: str, function: Callable) -> Any:
    """Run the indicated function inside a group with the given name.

    :param name: the name of the group
    :param function: the code to run

    :return: whatever the function returns

    Usage:

    .. code-block:: python

       load("group", "group")
       load("requests", "get")

       def default(_):
            group("visit product listing page", lambda: ...)

            def _add():
                ...

            group("add several products to the shopping cart", _add)
    """
