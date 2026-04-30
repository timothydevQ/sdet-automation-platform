from __future__ import annotations

import json
from pathlib import Path


def parse(path: Path) -> dict:
    data = json.loads(path.read_text())
    metrics = data.get("metrics", {})

    perf = []
    for name, m in metrics.items():
        values = m.get("values", {})
        for k, v in values.items():
            perf.append({"scenario": "checkout", "metric": f"{name}.{k}", "value": float(v)})

    threshold_results = data.get("root_group", {}).get("thresholds", {}) or data.get("thresholds", {})
    failed = sum(1 for t in threshold_results.values() if t and not t.get("ok", True))
    passed = max(0, len(threshold_results) - failed)

    return {
        "suite": "performance",
        "total": len(threshold_results),
        "passed": passed,
        "failed": failed,
        "skipped": 0,
        "duration_ms": int(metrics.get("iteration_duration", {}).get("values", {}).get("avg", 0)),
        "cases": [],
        "performance": perf,
    }
