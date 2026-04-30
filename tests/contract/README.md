# Contract tests

OpenAPI spec lives in `openapi.yaml`. Schemathesis generates fuzz inputs against
each operation and checks the responses match the schema.

```
pip install -r requirements.txt
API_BASE=http://localhost:8080 pytest -v
```
