import uuid


def email(prefix: str = "u") -> str:
    return f"{prefix}{uuid.uuid4().hex[:10]}@test.local"


def admin_email() -> str:
    return f"a{uuid.uuid4().hex[:10]}@admin.local"


def password() -> str:
    return "Hunter22!"
