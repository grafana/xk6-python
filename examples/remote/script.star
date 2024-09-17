"""
Example k6 test script to demonstrate the use of remote modules.
"""

load("https://grafana.github.io/k6pylib/welcome.py", "hello")

def default(_):
    hello("Joe")
