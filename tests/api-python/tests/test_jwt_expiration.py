import time

import jwt as pyjwt
import pytest


@pytest.mark.regression
@pytest.mark.auth
@pytest.mark.negative
def test_expired_customer_token_is_rejected(http, new_user):
    """Forge an already-expired token and verify the auth service rejects it."""
    import os
    secret = os.getenv("JWT_SECRET", "dev-secret-do-not-use-in-prod")
    expired_token = pyjwt.encode(
        {
            "uid": 999999,
            "role": "customer",
            "iss": "sdet-auth",
            "iat": int(time.time()) - 7200,
            "exp": int(time.time()) - 3600,
        },
        secret,
        algorithm="HS256",
    )
    r = http.get("/cart", headers={"Authorization": f"Bearer {expired_token}"})
    assert r.status_code in (401, 403), (
        f"Expired token should be rejected, got {r.status_code}"
    )


@pytest.mark.regression
@pytest.mark.auth
@pytest.mark.negative
def test_admin_expired_token_skew_bug(http):
    """BUG: admin tokens are accepted up to 60s past expiry due to clock-skew
    tolerance in auth-service/jwt.go parse(). This test asserts the DESIRED
    behavior (expired = rejected). Currently FAILS, exposing the bug.
    Fix: remove the 60s admin skew exemption in jwt.go.
    """
    import os
    secret = os.getenv("JWT_SECRET", "dev-secret-do-not-use-in-prod")
    just_expired = pyjwt.encode(
        {
            "uid": 999999,
            "role": "admin",
            "iss": "sdet-auth",
            "iat": int(time.time()) - 90,
            "exp": int(time.time()) - 30,
        },
        secret,
        algorithm="HS256",
    )
    r = http.get("/admin/orders", headers={"Authorization": f"Bearer {just_expired}"})
    assert r.status_code in (401, 403), (
        f"Expired admin token should be rejected within 60s of expiry, got {r.status_code}. "
        "Bug: auth-service/jwt.go grants 60s skew grace to admin role."
    )


@pytest.mark.regression
@pytest.mark.auth
@pytest.mark.negative
def test_malformed_token_rejected(http):
    r = http.get("/cart", headers={"Authorization": "Bearer not.a.real.token"})
    assert r.status_code in (401, 403)


@pytest.mark.regression
@pytest.mark.auth
@pytest.mark.negative
def test_missing_token_rejected(http):
    r = http.get("/cart")
    assert r.status_code == 401
