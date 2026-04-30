import uuid

import pytest

from helpers.cart import add_to_cart, checkout


@pytest.mark.smoke
@pytest.mark.orders
def test_checkout_success(http, new_user):
    add_to_cart(http, new_user, "SKU-001", 1)
    r = checkout(http, new_user)
    assert r.status_code == 201
    body = r.json()
    assert body["status"] == "paid"
    assert body["total_cents"] > 0


@pytest.mark.regression
@pytest.mark.orders
def test_checkout_with_coupon(http, new_user):
    add_to_cart(http, new_user, "SKU-001", 1)
    r = checkout(http, new_user, coupon="WELCOME10")
    assert r.status_code == 201


@pytest.mark.regression
@pytest.mark.orders
@pytest.mark.negative
def test_empty_cart_rejected(http, new_user):
    r = checkout(http, new_user)
    assert r.status_code == 400


@pytest.mark.regression
@pytest.mark.orders
def test_checkout_idempotency(http, new_user):
    add_to_cart(http, new_user, "SKU-001", 1)
    key = uuid.uuid4().hex
    r1 = checkout(http, new_user, idem_key=key)
    assert r1.status_code == 201
    add_to_cart(http, new_user, "SKU-001", 1)
    r2 = checkout(http, new_user, idem_key=key)
    assert r2.status_code == 200
    assert r1.json()["id"] == r2.json()["id"]


@pytest.mark.regression
@pytest.mark.orders
@pytest.mark.negative
def test_payment_decline(http, new_user):
    add_to_cart(http, new_user, "SKU-001", 1)
    r = checkout(http, new_user, card="tok_decline_card")
    assert r.status_code == 402
