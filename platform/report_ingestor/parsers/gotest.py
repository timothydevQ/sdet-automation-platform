from __future__ import annotations

import json
from pathlib import Path
from collections import defaultdict


def parse(path: Path) -> dict:
    cases: dict[str, dict] = {}
    elapsed_total = 0.0

    for line in path.read_text().splitlines():
        if not line.strip():
            continue
        try:
            evt = json.loads(line)
        except json.JSONDecodeError:
            continue
        test = evt.get("Test")
        if not test:
            continue
        action = evt.get("Action")
        case = cases.setdefault(test, {
            "name": test,
            "classname": evt.get("Package"),
            "status": "passed",
            "duration_ms": 0,
        })
        if action == "fail":
            case["status"] = "failed"
            case["failure"] = {"message": "test failed", "type": "go_failure", "stack": ""}
        elif action == "skip":
            case["status"] = "skipped"
        elif action == "pass":
            case["duration_ms"] = int(evt.get("Elapsed", 0) * 1000)
            elapsed_total += evt.get("Elapsed", 0)

    cs = list(cases.values())
    return {
        "suite": "go-integration",
        "total": len(cs),
        "passed": sum(1 for c in cs if c["status"] == "passed"),
        "failed": sum(1 for c in cs if c["status"] == "failed"),
        "skipped": sum(1 for c in cs if c["status"] == "skipped"),
        "duration_ms": int(elapsed_total * 1000),
        "cases": cs,
    }
