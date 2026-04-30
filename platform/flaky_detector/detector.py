from __future__ import annotations

import os
import statistics
from dataclasses import dataclass
from pathlib import Path

import psycopg
import typer
from rich.console import Console

DSN = os.getenv("INGESTOR_DSN", "postgres://sdet:sdet@localhost:5432/sdet")
LOOKBACK_RUNS = int(os.getenv("FLAKE_LOOKBACK_RUNS", "30"))
QUARANTINE_THRESHOLD = float(os.getenv("FLAKE_QUARANTINE", "0.7"))
FLAKY_THRESHOLD = float(os.getenv("FLAKE_THRESHOLD", "0.3"))

app = typer.Typer(add_completion=False)
console = Console()


@dataclass
class TestStat:
    suite: str
    classname: str | None
    name: str
    samples: int
    failures: int
    pass_after_retry: int
    durations_ms: list[int]
    recent_streak: int

    @property
    def failure_rate(self) -> float:
        return self.failures / self.samples if self.samples else 0.0

    @property
    def pass_after_retry_rate(self) -> float:
        return self.pass_after_retry / max(self.samples, 1)

    @property
    def duration_variance(self) -> float:
        if len(self.durations_ms) < 2:
            return 0.0
        m = statistics.mean(self.durations_ms)
        if m == 0:
            return 0.0
        return statistics.stdev(self.durations_ms) / m

    def score(self) -> float:
        return (
            0.5 * self.failure_rate
            + 0.3 * self.pass_after_retry_rate
            + 0.1 * min(1.0, self.duration_variance)
            + 0.1 * min(1.0, self.recent_streak / 5.0)
        )

    def classify(self) -> str:
        if self.score() >= QUARANTINE_THRESHOLD:
            return "quarantined"
        if self.failure_rate >= 0.9:
            return "broken"
        if self.score() >= FLAKY_THRESHOLD:
            return "flaky"
        if self.duration_variance > 0.5 and self.failure_rate < 0.05:
            return "slow"
        return "stable"


def fetch_stats(conn) -> list[TestStat]:
    with conn.cursor() as cur:
        cur.execute(
            """
            WITH recent AS (
                SELECT id FROM test_runs ORDER BY started_at DESC LIMIT %s
            )
            SELECT suite, classname, name, status, duration_ms, retries
            FROM test_cases
            WHERE run_id IN (SELECT id FROM recent)
            ORDER BY name
            """,
            (LOOKBACK_RUNS,),
        )
        rows = cur.fetchall()

    grouped: dict[tuple, dict] = {}
    for suite, classname, name, status, duration, retries in rows:
        key = (suite, classname, name)
        g = grouped.setdefault(key, {
            "samples": 0,
            "failures": 0,
            "pass_after_retry": 0,
            "durations_ms": [],
            "recent_results": [],
        })
        g["samples"] += 1
        g["recent_results"].append(status)
        if status in ("failed", "error"):
            g["failures"] += 1
        if status == "passed" and (retries or 0) > 0:
            g["pass_after_retry"] += 1
        if duration:
            g["durations_ms"].append(duration)

    stats: list[TestStat] = []
    for (suite, classname, name), g in grouped.items():
        streak = 0
        for s in reversed(g["recent_results"]):
            if s in ("failed", "error"):
                streak += 1
            else:
                break
        stats.append(TestStat(
            suite=suite,
            classname=classname,
            name=name,
            samples=g["samples"],
            failures=g["failures"],
            pass_after_retry=g["pass_after_retry"],
            durations_ms=g["durations_ms"],
            recent_streak=streak,
        ))
    return stats


def upsert_scores(conn, stats: list[TestStat]):
    with conn.cursor() as cur:
        for s in stats:
            cur.execute(
                """
                INSERT INTO flake_scores
                  (suite, classname, name, score, classification, failure_rate,
                   pass_after_retry_rate, duration_variance, recent_failure_streak, samples)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (suite, classname, name) DO UPDATE SET
                  score = EXCLUDED.score,
                  classification = EXCLUDED.classification,
                  failure_rate = EXCLUDED.failure_rate,
                  pass_after_retry_rate = EXCLUDED.pass_after_retry_rate,
                  duration_variance = EXCLUDED.duration_variance,
                  recent_failure_streak = EXCLUDED.recent_failure_streak,
                  samples = EXCLUDED.samples,
                  computed_at = NOW()
                """,
                (
                    s.suite, s.classname, s.name, s.score(), s.classify(),
                    s.failure_rate, s.pass_after_retry_rate, s.duration_variance,
                    s.recent_streak, s.samples,
                ),
            )


def render_flaky_md(stats: list[TestStat]) -> str:
    flaky = [s for s in stats if s.classify() in ("flaky", "quarantined", "broken")]
    flaky.sort(key=lambda s: s.score(), reverse=True)

    lines = [
        "# Flaky tests",
        "",
        f"Lookback: last {LOOKBACK_RUNS} runs.",
        f"Total tests analyzed: **{len(stats)}**.",
        f"Flaky / quarantined / broken: **{len(flaky)}**.",
        "",
        "| Test | Suite | Class | Score | Failure rate | Streak | Samples |",
        "|------|-------|-------|------:|-------------:|-------:|--------:|",
    ]
    for s in flaky[:50]:
        lines.append(
            f"| {s.classify()} `{s.name}` | {s.suite} | {s.classname or ''} | "
            f"{s.score():.2f} | {s.failure_rate*100:.0f}% | {s.recent_streak} | {s.samples} |"
        )
    if not flaky:
        lines.append("")
        lines.append("_No flaky tests detected over the lookback window._")
    return "\n".join(lines) + "\n"


def render_release_readiness(stats: list[TestStat]) -> tuple[str, bool]:
    by_class = {"stable": 0, "slow": 0, "flaky": 0, "broken": 0, "quarantined": 0}
    for s in stats:
        by_class[s.classify()] = by_class.get(s.classify(), 0) + 1

    blocking = by_class["broken"] > 0 or by_class["flaky"] > 5
    verdict = "GO" if not blocking else "NO-GO"

    lines = [
        "# Release readiness",
        "",
        f"**Verdict: {verdict}**",
        "",
        "| Classification | Count |",
        "|----------------|------:|",
    ]
    for k, v in by_class.items():
        lines.append(f"| {k} | {v} |")

    lines += ["", "## Reasoning", ""]
    if by_class["broken"]:
        lines.append(f"- {by_class['broken']} broken test(s) (>=90% failure rate). Block release.")
    if by_class["flaky"] > 5:
        lines.append(f"- {by_class['flaky']} flaky test(s) above threshold of 5. Investigate before release.")
    if not blocking:
        lines.append("- All checks pass. Suite is healthy enough to ship.")
    return "\n".join(lines) + "\n", blocking


@app.command()
def main(
    out: Path = typer.Option(Path("flaky-tests.md"), "--out"),
    release_check: bool = typer.Option(False, "--release-check"),
    dsn: str = typer.Option(DSN, "--dsn"),
):
    with psycopg.connect(dsn, autocommit=True) as conn:
        stats = fetch_stats(conn)
        upsert_scores(conn, stats)

    if release_check:
        md, blocking = render_release_readiness(stats)
        out.write_text(md)
        console.print(f"[bold]{'NO-GO' if blocking else 'GO'}[/] -> {out}")
        raise typer.Exit(code=1 if blocking else 0)

    out.write_text(render_flaky_md(stats))
    console.print(f"[green]wrote[/] {out} ({len(stats)} tests analyzed)")


if __name__ == "__main__":
    app()
