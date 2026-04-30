from __future__ import annotations

from pathlib import Path

from lxml import etree


def parse(path: Path) -> dict:
    root = etree.parse(str(path)).getroot()
    suites = root.findall(".//testsuite") or [root]

    total = passed = failed = skipped = 0
    duration_ms = 0
    cases: list[dict] = []

    for s in suites:
        for tc in s.findall("testcase"):
            total += 1
            name = tc.get("name", "?")
            classname = tc.get("classname")
            time_s = float(tc.get("time", 0) or 0)
            duration_ms += int(time_s * 1000)
            if tc.find("failure") is not None or tc.find("error") is not None:
                failed += 1
                node = tc.find("failure") if tc.find("failure") is not None else tc.find("error")
                cases.append({
                    "name": name,
                    "classname": classname,
                    "status": "failed" if node.tag == "failure" else "error",
                    "duration_ms": int(time_s * 1000),
                    "failure": {
                        "message": node.get("message"),
                        "type": node.get("type"),
                        "stack": (node.text or "")[:8000],
                    },
                })
            elif tc.find("skipped") is not None:
                skipped += 1
                cases.append({"name": name, "classname": classname, "status": "skipped"})
            else:
                passed += 1
                cases.append({
                    "name": name,
                    "classname": classname,
                    "status": "passed",
                    "duration_ms": int(time_s * 1000),
                })

    suite = root.get("name") or (suites[0].get("name") if suites else "unknown")
    return {
        "suite": suite,
        "total": total,
        "passed": passed,
        "failed": failed,
        "skipped": skipped,
        "duration_ms": duration_ms,
        "cases": cases,
    }
