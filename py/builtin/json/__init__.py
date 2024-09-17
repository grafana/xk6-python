"""
The json module helps with from python objects to json encoded strings and back
"""


def loads(s: str) -> object:
    """Loads the given JSON encoded string into a python object

    The function will throw an error if the string is not valid json

    (note: this function is also available under the alias 'decode')

    :param s: the string to decode

    :return: the decoded string s as a python object
    """

def dumps(obj: object) -> str:
    """Dumps the given python object into a json encoded string

    (note: this function is also available under the alias 'encode')

    :param obj: the python object to encode

    :return: an encoded string with the contents of the python object
    """
