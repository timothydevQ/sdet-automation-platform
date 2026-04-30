import pytest
from factories import users


@pytest.mark.smoke
@pytest.mark.auth
def test_login_success(http, new_user):
    r = http.post("/auth/login", json={"email": new_user.email, "password": new_user.password})
    assert r.status_code == 200
    assert "token" in r.json()


@pytest.mark.regression
@pytest.mark.auth
@pytest.mark.negative
def test_login_invalid_password(http, new_user):
    r = http.post("/auth/login", json={"email": new_user.email, "password": "wrong"})
    assert r.status_code == 401


@pytest.mark.regression
@pytest.mark.auth
@pytest.mark.negative
def test_login_unknown_email(http):
    r = http.post("/auth/login", json={"email": users.email(), "password": "whatever"})
    assert r.status_code == 401


@pytest.mark.regression
@pytest.mark.auth
def test_register_duplicate_email(http, new_user):
    r = http.post("/auth/register", json={"email": new_user.email, "password": "AnotherPass1!"})
    assert r.status_code == 409
