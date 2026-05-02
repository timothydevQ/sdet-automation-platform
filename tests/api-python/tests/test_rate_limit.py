import pytest


@pytest.mark.regression
@pytest.mark.negative
def test_rate_limit_enforced(http):
    seen_429 = False
    for _ in range(150):
        r = http.get("/healthz")
        if r.status_code == 429:
            seen_429 = True
            break
    assert seen_429, "Rate limiter should block after threshold"


@pytest.mark.regression
@pytest.mark.negative
def test_rate_limit_bypass_via_xff(http):
    """BUG: X-Forwarded-For rotation bypasses the rate limiter.
    This test asserts the DESIRED behavior (bypass should not work).
    It currently FAILS, exposing the bug in api-gateway/main.go clientKey().
    Fix: validate and ignore untrusted X-Forwarded-For headers.
    """
    blocked = 0
    for i in range(150):
        r = http.get("/healthz", headers={"X-Forwarded-For": f"10.0.0.{i % 250}"})
        if r.status_code == 429:
            blocked += 1
    assert blocked > 0, (
        "Rate limiter must not be bypassable via X-Forwarded-For rotation. "
        "Bug in api-gateway/main.go clientKey() — trusts XFF unconditionally."
    )


@pytest.mark.regression
@pytest.mark.negative
def test_admin_role_check_is_case_sensitive(http, new_user, auth_headers):
    """Admin role check is case-sensitive; tokens with role 'Admin' (capital A)
    are incorrectly rejected. This test asserts the customer path is blocked."""
    r = http.get("/admin/orders", headers=auth_headers(new_user))
    assert r.status_code == 403, f"Customer should be forbidden, got {r.status_code}"
