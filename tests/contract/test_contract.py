import os

import schemathesis

API_BASE = os.getenv("API_BASE", "http://localhost:8080")
schema = schemathesis.from_path("openapi.yaml", base_url=API_BASE)


@schema.parametrize()
def test_api_conforms_to_schema(case):
    case.call_and_validate()
