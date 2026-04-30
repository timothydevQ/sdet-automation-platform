import os
import time
import uuid
from dataclasses import dataclass

import httpx
import pytest
from faker import Faker

BASE = os.getenv("API_BASE", "http://localhost:8080")
faker = Faker()


@dataclass
class User:
    email: str
    password: str
    token: str
    role: str


@pytest.fixture(scope="session")
def base_url():
    return BASE


@pytest.fixture(scope="session")
def http():
    with httpx.Client(base_url=BASE, timeout=10.0) as c:
        for _ in range(60):
            try:
                r = c.get("/healthz")
                if r.status_code == 200:
                    break
            except Exception:
                pass
            time.sleep(1)
        yield c


@pytest.fixture
def new_user(http):
    email = f"u{uuid.uuid4().hex[:10]}@test.local"
    password = "Hunter22!"
    r = http.post("/auth/register", json={"email": email, "password": password})
    assert r.status_code == 201, r.text
    body = r.json()
    return User(email=email, password=password, token=body["token"], role=body["role"])


@pytest.fixture
def admin_user(http):
    email = f"a{uuid.uuid4().hex[:10]}@admin.local"
    password = "Hunter22!"
    r = http.post("/auth/register", json={"email": email, "password": password})
    assert r.status_code == 201
    body = r.json()
    return User(email=email, password=password, token=body["token"], role=body["role"])


@pytest.fixture
def auth_headers():
    def make(user: User, extra: dict | None = None) -> dict:
        h = {"Authorization": f"Bearer {user.token}"}
        if extra:
            h.update(extra)
        return h

    return make
