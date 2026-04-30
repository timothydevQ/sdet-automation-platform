from __future__ import annotations

import json
from pathlib import Path

from . import junit as junit_parser


def parse(path: Path) -> dict:
    if path.suffix == ".xml":
        return junit_parser.parse(path)

    data = json.loads(path.read_text())
    cases: list[dict] = []
    total = passed = failed = skipped = 0
    duration_ms = 0

    def walk(items):
        nonlocal total, passed, failed, skipped, duration_ms
        for it in items:
            for spec in it.get("specs", []):
                for t in spec.get("tests", []):
                    total += 1
                    status = "passed"
                    last = t["results"][-1] if t.get("results") else {}
                    s = last.get("status", "passed")
                    if s == "passed":
                        passed += 1
                    elif s == "skipped":
                        skipped += 1
                        status = "skipped"
                    else:
                        failed += 1
                        status = "failed"
                    duration_ms += int(last.get("duration", 0))
                    cases.append({
                        "name": spec.get("title"),
                        "classname": it.get("title"),
                        "status": status,
                        "duration_ms": int(last.get("duration", 0)),
                        "retries": len(t.get("results", [])) - 1,
                    })
            walk(it.get("suites", []))

    walk(data.get("suites", []))
    return {
        "suite": "playwright",
        "total": total,
        "passed": passed,
        "failed": failed,
        "skipped": skipped,
        "duration_ms": duration_ms,
        "cases": cases,
    }
