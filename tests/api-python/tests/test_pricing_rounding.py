import pytest

from helpers.cart import add_to_cart, checkout


@pytest.mark.regression
@pytest.mark.orders
def test_welcome10_exact_discount(http, new_user):
    """WELCOME10 removes exactly 10% — verify to the cent."""
    add_to_cart(http, new_user, "SKU-001", 1)
    r = checkout(http, new_user, coupon="WELCOME10")
    assert r.status_code == 201
    body = r.json()
    unit = 19999
    expected_total = unit - unit // 10
    assert body["total_cents"] == expected_total, (
        f"WELCOME10 should yield {expected_total} cents, got {body['total_cents']}"
    )


@pytest.mark.regression
@pytest.mark.orders
def test_bulk20_rounding_bug(http, new_user):
    """BUG: BULK20 uses integer division which rounds wrong at certain amounts.
    This test asserts the DESIRED behavior (mathematically correct discount).
    It currently FAILS for some subtotals, exposing the bug in order-service/orders.go.
    Fix: use proper rounding in applyDiscount().
    """
    add_to_cart(http, new_user, "SKU-001", 2)
    r = checkout(http, new_user, coupon="BULK20")
    assert r.status_code == 201
    body = r.json()
    subtotal = 19999 * 2
    expected_discount = round(subtotal * 0.20)
    expected_total = subtotal - expected_discount
    assert body["total_cents"] == expected_total, (
        f"BULK20 should discount exactly 20% (rounded). "
        f"Expected {expected_total}, got {body['total_cents']}. "
        f"Bug in order-service/orders.go applyDiscount()."
    )


@pytest.mark.regression
@pytest.mark.orders
@pytest.mark.negative
def test_invalid_coupon_no_discount(http, new_user):
    add_to_cart(http, new_user, "SKU-001", 1)
    r = checkout(http, new_user, coupon="FAKECOUPON")
    assert r.status_code == 201
    body = r.json()
    assert body["total_cents"] == 19999, (
        f"Unknown coupon should apply no discount, got {body['total_cents']}"
    )
