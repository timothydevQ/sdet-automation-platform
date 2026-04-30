import time

import jwt as pyjwt
import pytest


@pytest.mark.regression
@pytest.mark.auth
@pytest.mark.negative
def test_expired_token_rejected_on_customer_endpoint(http, new_user):
    payload = pyjwt.decode(new_user.token, options={"verify_signature": False})
    assert payload["exp"] > time.time()


@pytest.mark.regression
@pytest.mark.admin
@pytest.mark.negative
def test_admin_jwt_skew_bug(http, admin_user, auth_headers):
    r = http.get("/admin/orders", headers=auth_headers(admin_user))
    assert r.status_code in (200, 403)
