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
    assert seen_429


@pytest.mark.regression
@pytest.mark.negative
def test_rate_limit_bypass_via_xff(http):
    """X-Forwarded-For rotation currently bypasses the limiter (known bug)."""
    blocked = 0
    for i in range(150):
        r = http.get("/healthz", headers={"X-Forwarded-For": f"10.0.0.{i % 250}"})
        if r.status_code == 429:
            blocked += 1
    assert blocked == 0


@pytest.mark.regression
@pytest.mark.negative
def test_admin_authorization_case(http, new_user, auth_headers):
    r = http.get("/admin/orders", headers=auth_headers(new_user))
    assert r.status_code == 403
