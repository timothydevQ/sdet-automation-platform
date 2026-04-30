from __future__ import annotations

import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Iterable

import psycopg
import typer
from rich.console import Console

from .parsers import junit, playwright as pw, k6, nunit, gotest

app = typer.Typer(add_completion=False)
console = Console()

DEFAULT_DSN = os.getenv(
    "INGESTOR_DSN",
    "postgres://sdet:sdet@localhost:5432/sdet",
)


def discover(root: Path) -> Iterable[tuple[str, Path]]:
    if not root.exists():
        return []
    for p in root.rglob("*"):
        if not p.is_file():
            continue
        name = p.name.lower()
        if name.endswith(".xml") and "junit" in name:
            yield "junit", p
        elif name.endswith(".xml") and "surefire" in p.parts[-2].lower() if len(p.parts) >= 2 else False:
            yield "junit", p
        elif name == "playwright-report.json" or "playwright-junit" in name:
            yield "playwright", p
        elif name.startswith("k6-") and name.endswith(".json"):
            yield "k6", p
        elif "nunit" in name and name.endswith(".xml"):
            yield "nunit", p
        elif name == "go-test.json":
            yield "gotest", p


def insert_run(conn, suite: str, parsed: dict) -> int:
    with conn.cursor() as cur:
        cur.execute(
            """
            INSERT INTO test_runs (suite, branch, commit_sha, workflow_run,
                started_at, finished_at, total, passed, failed, skipped,
                duration_ms, metadata)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s::jsonb)
            RETURNING id
            """,
            (
                suite,
                os.getenv("GITHUB_REF_NAME"),
                os.getenv("GITHUB_SHA"),
                os.getenv("GITHUB_RUN_ID"),
                parsed.get("started_at") or datetime.now(timezone.utc),
                parsed.get("finished_at") or datetime.now(timezone.utc),
                parsed["total"],
                parsed["passed"],
                parsed["failed"],
                parsed["skipped"],
                parsed.get("duration_ms"),
                json.dumps(parsed.get("metadata", {})),
            ),
        )
        return cur.fetchone()[0]


def insert_cases(conn, run_id: int, suite: str, cases: list[dict]):
    if not cases:
        return
    with conn.cursor() as cur:
        cur.executemany(
            """
            INSERT INTO test_cases (run_id, suite, classname, name, status,
                duration_ms, retries)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
            RETURNING id
            """,
            [(run_id, suite, c.get("classname"), c["name"], c["status"],
              c.get("duration_ms"), c.get("retries", 0)) for c in cases],
            returning=True,
        )
        ids = []
        while True:
            rec = cur.fetchone()
            if rec is None:
                break
            ids.append(rec[0])
            if not cur.nextset():
                break

    failures = []
    for c, cid in zip(cases, ids):
        if c["status"] in ("failed", "error") and c.get("failure"):
            failures.append((cid, c["failure"].get("message"),
                             c["failure"].get("stack"),
                             c["failure"].get("type")))
    if failures:
        with conn.cursor() as cur:
            cur.executemany(
                "INSERT INTO test_failures (case_id, message, stack, failure_type) VALUES (%s, %s, %s, %s)",
                failures,
            )


def ingest_one(conn, kind: str, path: Path):
    parsers = {
        "junit": junit.parse,
        "playwright": pw.parse,
        "k6": k6.parse,
        "nunit": nunit.parse,
        "gotest": gotest.parse,
    }
    parser = parsers.get(kind)
    if parser is None:
        console.print(f"[yellow]skip {path}: unknown kind {kind}[/]")
        return
    try:
        parsed = parser(path)
    except Exception as e:
        console.print(f"[red]parse error {path}: {e}[/]")
        return
    suite = parsed.get("suite") or kind
    run_id = insert_run(conn, suite, parsed)
    insert_cases(conn, run_id, suite, parsed.get("cases", []))
    if parsed.get("performance"):
        with conn.cursor() as cur:
            cur.executemany(
                """
                INSERT INTO performance_results (run_id, scenario, metric, value, threshold, passed)
                VALUES (%s, %s, %s, %s, %s, %s)
                """,
                [(run_id, p["scenario"], p["metric"], p["value"], p.get("threshold"), p.get("passed"))
                 for p in parsed["performance"]],
            )
    console.print(f"[green]ingested[/] {kind} {path.name} ({parsed['total']} tests)")


@app.command()
def main(artifacts_dir: Path = typer.Argument(..., exists=True),
         dsn: str = typer.Option(DEFAULT_DSN, "--dsn")):
    with psycopg.connect(dsn, autocommit=True) as conn:
        for kind, path in discover(artifacts_dir):
            ingest_one(conn, kind, path)


if __name__ == "__main__":
    app()
