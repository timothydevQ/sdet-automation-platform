import uuid


def cart_payload(sku: str = "SKU-001", qty: int = 1) -> dict:
    return {"sku": sku, "qty": qty}


def checkout_payload(coupon: str = "", card: str = "tok_test_visa") -> dict:
    return {"coupon": coupon, "card_token": card}


def idempotency_key() -> str:
    return f"idem-{uuid.uuid4().hex}"
