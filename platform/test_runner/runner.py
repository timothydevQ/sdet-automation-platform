from __future__ import annotations

import subprocess
from pathlib import Path

import typer
from rich.console import Console

app = typer.Typer(add_completion=False)
console = Console()

ROOT = Path(__file__).resolve().parents[2]

SUITES = {
    "api": ["make", "-C", str(ROOT), "test-api"],
    "e2e": ["make", "-C", str(ROOT), "test-e2e"],
    "java": ["make", "-C", str(ROOT), "test-java"],
    "dotnet": ["make", "-C", str(ROOT), "test-dotnet"],
    "ruby": ["make", "-C", str(ROOT), "test-ruby"],
    "go": ["make", "-C", str(ROOT), "test-go"],
    "cpp": ["make", "-C", str(ROOT), "test-cpp"],
    "perf": ["make", "-C", str(ROOT), "test-performance"],
    "security": ["make", "-C", str(ROOT), "test-security"],
}


@app.command()
def run(suite: str = typer.Argument(...)):
    if suite == "all":
        for name, cmd in SUITES.items():
            console.rule(f"[bold]{name}")
            subprocess.run(cmd, check=False)
        return
    cmd = SUITES.get(suite)
    if not cmd:
        console.print(f"[red]unknown suite: {suite}[/]")
        raise typer.Exit(2)
    raise SystemExit(subprocess.run(cmd).returncode)


if __name__ == "__main__":
    app()
