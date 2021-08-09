import pytest

from server.app import create_app


@pytest.fixture
def client():
    app = create_app({'TESTING': True})


def test_checkout(client, number, exp_month, exp_year, cvc):
    resp = client.post("/checkout", data=dict(
        number=number,
        exp_month=exp_month,
        exp_year=exp_year,
        cvc=cvc
    ), follow_redirects=True)
    assert b"'status': 'success'" in resp.data
