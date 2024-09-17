"""
The welcome module is an example embedded module.
The file extension can be .py or .star, the advantage of .star is that
using the bazel-stack-vscode Visual Studio Code extension,
quite good editor support is available for the starlark language.
"""

def hello(name = "World"):
    """
    The hello function is an example function implemented in starlark.
    """
    print("Hello, %s!" % (name))  # buildifier: disable=print
