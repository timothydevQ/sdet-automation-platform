def add_to_cart(http, user, sku: str, qty: int = 1):
    r = http.post(
        "/cart/items",
        json={"sku": sku, "qty": qty},
        headers={"Authorization": f"Bearer {user.token}"},
    )
    r.raise_for_status()
    return r.json()


def checkout(http, user, coupon: str = "", card: str = "tok_test_visa", idem_key: str | None = None):
    headers = {"Authorization": f"Bearer {user.token}"}
    if idem_key:
        headers["Idempotency-Key"] = idem_key
    return http.post(
        "/checkout",
        json={"coupon": coupon, "card_token": card},
        headers=headers,
    )
