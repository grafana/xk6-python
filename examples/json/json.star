"""
Example k6 test script to demonstrate the use encoding and decoding json.
"""

load("json", "dumps", "loads")

def default(_):
    raw_data = '{"field":"string","another_field":[1,2]}'
    print(raw_data)
    data = loads(raw_data)
    print(data)
    encoded = dumps(data)
    print(encoded)
