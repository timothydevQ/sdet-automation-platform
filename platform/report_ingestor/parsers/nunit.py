from __future__ import annotations

from pathlib import Path

from . import junit as junit_parser


def parse(path: Path) -> dict:
    out = junit_parser.parse(path)
    out["suite"] = "nunit"
    return out
