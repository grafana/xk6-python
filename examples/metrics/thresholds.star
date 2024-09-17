"""
Example k6 test script to demonstrate the use custom metrics and thresholds.
"""

load("requests", "get")
load("metrics", "counter")

errorCounter = counter("errors")

options = { 'thresholds': { 'errors': ['count<100'] } }

def default(_):
    """
    The default function defines the logic of the VUs.
    """
    resp = get('https://test-api.k6.io/public/crocodiles/1/')

    contentOK = resp.json()['name'] == 'Bert'
    errorCounter.add(not contentOK)
