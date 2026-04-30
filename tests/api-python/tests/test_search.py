import pytest


@pytest.mark.regression
@pytest.mark.catalog
def test_search_basic(http):
    r = http.get("/catalog/products?q=widget")
    assert r.status_code == 200


@pytest.mark.regression
@pytest.mark.catalog
@pytest.mark.negative
@pytest.mark.parametrize("q", ["50%", "_underscore", "'quote", "100%off"])
def test_search_special_chars(http, q):
    r = http.get(f"/catalog/products?q={q}")
    assert r.status_code in (200, 400)
    if r.status_code == 200:
        r.json()
