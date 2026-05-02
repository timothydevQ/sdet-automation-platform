import threading
import pytest

from helpers.cart import add_to_cart, checkout


@pytest.mark.regression
@pytest.mark.catalog
@pytest.mark.negative
def test_out_of_stock_blocks_checkout(http, new_user):
    add_to_cart(http, new_user, "SKU-LIMITED", 999)
    r = checkout(http, new_user)
    assert r.status_code == 409, f"Expected 409 out-of-stock, got {r.status_code}"


@pytest.mark.regression
@pytest.mark.catalog
def test_successful_checkout_decrements_stock(http, new_user):
    """Stock should decrease by exactly the quantity purchased."""
    before = http.get("/catalog/products?q=Keyboard").json()
    stock_before = next((p["stock"] for p in before if p["sku"] == "SKU-001"), None)
    assert stock_before is not None

    add_to_cart(http, new_user, "SKU-001", 1)
    r = checkout(http, new_user)
    assert r.status_code == 201

    after = http.get("/catalog/products?q=Keyboard").json()
    stock_after = next((p["stock"] for p in after if p["sku"] == "SKU-001"), None)
    assert stock_after == stock_before - 1, (
        f"Stock should decrease by 1: {stock_before} -> {stock_after}"
    )


@pytest.mark.regression
@pytest.mark.catalog
@pytest.mark.flaky
def test_concurrent_checkout_inventory_race(http, new_user, admin_user):
    """BUG: non-atomic stock decrement in catalog-service/main.go reserveStock().
    Two concurrent checkouts against a single-unit item can both succeed,
    driving stock negative. This test asserts the DESIRED invariant.
    Fix: use SELECT ... FOR UPDATE in the reserve query.
    """
    add_to_cart(http, new_user, "SKU-RACE", 1)
    add_to_cart(http, admin_user, "SKU-RACE", 1)

    results = []
    def go(user):
        results.append(checkout(http, user))

    t1 = threading.Thread(target=go, args=(new_user,))
    t2 = threading.Thread(target=go, args=(admin_user,))
    t1.start(); t2.start(); t1.join(); t2.join()

    statuses = sorted(r.status_code for r in results)
    assert statuses == [201, 409], (
        f"Exactly one checkout should succeed and one should fail with 409. "
        f"Got statuses: {statuses}. "
        f"Bug in catalog-service/main.go reserveStock() — non-atomic decrement."
    )

    stock = http.get("/catalog/products?q=RACE").json()
    final_stock = next((p["stock"] for p in stock if p["sku"] == "SKU-RACE"), -1)
    assert final_stock >= 0, (
        f"Stock went negative ({final_stock}) — race condition confirmed."
    )
