import pytest

from helpers.cart import add_to_cart, checkout


@pytest.mark.regression
@pytest.mark.orders
@pytest.mark.parametrize("sku,qty,expected_min", [
    ("SKU-001", 5, 1),
    ("SKU-002", 3, 1),
])
def test_bulk_discount_rounding(http, new_user, sku, qty, expected_min):
    add_to_cart(http, new_user, sku, qty)
    r = checkout(http, new_user, coupon="BULK20")
    assert r.status_code == 201
    body = r.json()
    assert body["total_cents"] >= expected_min
