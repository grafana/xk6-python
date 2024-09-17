"""
The requests module allows issuing HTTP requests in k6 tests.
The usual k6 metrics of requests will be generated.
"""


def request(method, url, **kwargs):
    """Constructs and sends a HTTP request.

    :param method: method for the HTTP request: ``GET``, ``OPTIONS``, ``HEAD``, ``POST``, ``PUT``, ``PATCH``, or ``DELETE``.
    :param url: URL for the HTTP request.
    :param params: (optional) Dictionary, list of tuples or bytes to send
        in the query string for the HTTP request.
    :param data: (optional) Dictionary, list of tuples, bytes or string
        to send in the body of the HTTP request.
    :param json: (optional) A JSON serializable Python object to send in the body of the HTTP request.
    :param headers: (optional) Dictionary of HTTP Headers to send with the HTTP request.
    :param cookies: (optional) Dict or CookieJar object to send with the HTTP request.
    :return: `Response` struct
    :rtype: Response

    Usage:

    .. code-block:: python

       load("requests", "request")

       def default(_):
           resp = request("GET", "https://httpbin.org/get")
    """


def get(url, params=None, **kwargs):
    """Sends a GET request.

    :param url: URL for the HTTP request.
    :param params: (optional) Dictionary, list of tuples or bytes to send
        in the query string for the HTTP request.
    :param kwargs: Optional keyword arguments that ``request`` takes.
    :return: `Response` struct
    :rtype: Response

    Usage:

    .. code-block:: python

       load("requests", "get")

       def default(_):
           resp = get("https://httpbin.org/get")
    """


def options(url, **kwargs):
    """Sends an OPTIONS request.

    :param url: URL for the HTTP request.
    :param kwargs: Optional keyword arguments that ``request`` takes.
    :return: `Response` struct
    :rtype: Response

    Usage:

    .. code-block:: python

       load("requests", "options")

       def default(_):
           resp = options("https://httpbin.org/options")
    """


def head(url, **kwargs):
    """Sends a HEAD request.

    :param url: URL for the new HTTP request.
    :param kwargs: Optional keyword arguments that ``request`` takes.
    :return: `Response` struct
    :rtype: Response

    Usage:

    .. code-block:: python

       load("requests", "head")

       def default(_):
           head = head("https://httpbin.org/head")
    """


def post(url, data=None, json=None, **kwargs):
    """Sends a POST request.

    :param url: URL for the new HTTP request.
    :param data: (optional) Dictionary, list of tuples, bytes, or file-like
        object to send in the body of the HTTP request.
    :param json: (optional) A JSON serializable Python object to send in the body of the HTTP request.
    :param kwargs: Optional keyword arguments that ``request`` takes.
    :return: `Response` struct
    :rtype: Response

    Usage:

    .. code-block:: python

       load("requests", "post")

       def default(_):
           resp = post("https://httpbin.org/post")
    """


def put(url, data=None, **kwargs):
    """Sends a PUT request.

    :param url: URL for the new HTTP request.
    :param data: (optional) Dictionary, list of tuples, bytes, or file-like
        object to send in the body of the HTTP request.
    :param json: (optional) A JSON serializable Python object to send in the body of the HTTP request.
    :param kwargs: Optional keyword arguments that ``request`` takes.
    :return: `Response` struct
    :rtype: Response

    Usage:

    .. code-block:: python

       load("requests", "put")

       def default(_):
           resp = put("https://httpbin.org/put")
    """


def patch(url, data=None, **kwargs):
    """Sends a PATCH request.

    :param url: URL for the new HTTP request.
    :param data: (optional) Dictionary, list of tuples, bytes, or file-like
        object to send in the body of the HTTP request.
    :param json: (optional) A JSON serializable Python object to send in the body of the HTTP request.
    :param kwargs: Optional keyword arguments that ``request`` takes.
    :return: `Response` struct
    :rtype: Response

    Usage:

    .. code-block:: python

       load("requests", "patch")

       def default(_):
           resp = patch("https://httpbin.org/patch")
    """


def delete(url, **kwargs):
    """Sends a DELETE request.

    :param url: URL for the new HTTP request.
    :param kwargs: Optional keyword arguments that ``request`` takes.
    :return: :class:`Response` struct
    :rtype: Response

    Usage:

    .. code-block:: python

       load("requests", "delete")

       def default(_):
           resp = delete("https://httpbin.org/delete")
    """


class Response:
    """The `Response` struct, which contains a server's response to an HTTP request.
    """

    @property
    def status_code(self):
        """Integer Code of responded HTTP Status, e.g. 404 or 200."""

    @property
    def headers(self):
        """Dictionary of Response Headers.
        For example, ``headers['content-encoding']`` will return the
        value of a ``'Content-Encoding'`` response header.
        """

    @property
    def url(self):
        """Final URL location of Response."""

    @property
    def reason(self):
        """Textual reason of responded HTTP Status, e.g. "Not Found" or "OK"."""

    @property
    def ok(self):
        """Returns True if :attr:`status_code` is less than 400, False if not.

        This attribute checks if the status code of the response is between
        400 and 600 to see if there was a client error or a server error. If
        the status code is between 200 and 400, this will return True. This
        is **not** a check to see if the response code is ``200 OK``.
        """

    @property
    def text(self):
        """Content of the response, in unicode."""

    def json(self, **kwargs):
        """Returns the json-encoded content of a response, if any.

        :param kwargs: Optional keyword arguments that ``json.loads`` takes.
        """
