import pytest


@pytest.mark.smoke
@pytest.mark.admin
@pytest.mark.negative
def test_customer_cannot_list_admin_orders(http, new_user, auth_headers):
    r = http.get("/admin/orders", headers=auth_headers(new_user))
    assert r.status_code == 403


@pytest.mark.regression
@pytest.mark.admin
def test_admin_can_list_orders(http, admin_user, auth_headers):
    r = http.get("/admin/orders", headers=auth_headers(admin_user))
    assert r.status_code == 200
    assert isinstance(r.json(), list)


@pytest.mark.regression
@pytest.mark.admin
@pytest.mark.negative
def test_no_token_rejected(http):
    r = http.get("/admin/orders")
    assert r.status_code == 401
