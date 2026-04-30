import uuid


def customer_email() -> str:
    return f"u{uuid.uuid4().hex[:10]}@test.local"


def admin_email() -> str:
    return f"a{uuid.uuid4().hex[:10]}@admin.local"


def password() -> str:
    return "Hunter22!"


def credentials(role: str = "customer") -> tuple[str, str]:
    return (admin_email() if role == "admin" else customer_email(), password())
