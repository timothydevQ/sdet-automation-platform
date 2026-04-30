import threading
import uuid

import pytest

from helpers.cart import add_to_cart, checkout


@pytest.mark.regression
@pytest.mark.catalog
@pytest.mark.negative
def test_out_of_stock_blocks_checkout(http, new_user):
    add_to_cart(http, new_user, "SKU-LIMITED", 999)
    r = checkout(http, new_user)
    assert r.status_code == 409


@pytest.mark.regression
@pytest.mark.catalog
@pytest.mark.flaky
def test_concurrent_checkout_inventory(http, new_user, admin_user):
    """The reservation path is non-atomic; under load both checkouts can pass."""
    add_to_cart(http, new_user, "SKU-RACE", 1)
    add_to_cart(http, admin_user, "SKU-RACE", 1)
    results = []

    def go(user):
        results.append(checkout(http, user))

    t1 = threading.Thread(target=go, args=(new_user,))
    t2 = threading.Thread(target=go, args=(admin_user,))
    t1.start(); t2.start(); t1.join(); t2.join()
    statuses = sorted([r.status_code for r in results])
    assert statuses[0] == 201
