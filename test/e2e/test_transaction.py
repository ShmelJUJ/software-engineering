from pytest import fixture
from time import sleep
import json
from urllib.request import urlopen, Request

TRANSACTION_URL = "http://localhost:8083"
API_PATH = "/api/v1"
CREATE_TRANSACTION_PATH = "/transaction/create"


DATA_TO_CREATE_TRANSACTION = {
        "idempotency_token": "1234567890",
        "money_info": {
            "method": "bank_account",
            "currency": "USD",
            "amount": 100
        },
        "receiver": {
            "user_id": "00000000-0000-0000-0000-000000000000"
        },
        "sender": {
            "user_id": "11111111-1111-1111-1111-111111111111"
        }
    }

DATA_TO_ACCEPT_TRANSACTION = {
    "sender": {
        "user_id": "11111111-1111-1111-1111-111111111111"
    }
}



def headers():
    return {'content-type': 'application/json'}


def create_transaction_send(data: list) -> dict:
    req = Request(TRANSACTION_URL + API_PATH + CREATE_TRANSACTION_PATH, data=json.dumps(
        data).encode(), headers=headers(), method='POST')
    response = urlopen(req)
    assert response.getcode() == 200
    result = json.loads(response.read().decode())
    return result

def status_transaction_send(id: str) -> list:
    req = Request(TRANSACTION_URL + API_PATH + "/transaction/" + id + "/retrieve/status", headers=headers(),
                  data=json.dumps(dict({"transaction_id": id})).encode(), method='GET')
    response = urlopen(req)
    assert response.getcode() == 200
    return json.loads(response.read().decode())

def cancel_transaction_send(id: str):
    req = Request(TRANSACTION_URL + API_PATH + "/transaction/" + id + "/cancel", headers=headers(),
                  data=json.dumps(dict({"transaction_id": id})).encode(), method='POST')
    response = urlopen(req)
    assert response.getcode() == 200

def accept_transaction_send(id: str, data: dict):
    req = Request(TRANSACTION_URL + API_PATH + "/transaction/" + id + "/accept", headers=headers(),
                  data=json.dumps(data).encode(), method='POST')
    response = urlopen(req)
    assert response.getcode() == 200

def retrieve_transaction_send(id: str):
    req = Request(TRANSACTION_URL + API_PATH + "/transaction/" + id + "/retrieve", headers=headers(), method='GET')
    response = urlopen(req)
    assert response.getcode() == 200
    return json.loads(response.read().decode())


def test_creation_transaction():
    result = create_transaction_send(DATA_TO_CREATE_TRANSACTION)

    id = result["transaction_id"]

    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = status_transaction_send(id)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0
    assert status_result["transaction_status"] == "created"

def test_cancel_transaction():
    result = create_transaction_send(DATA_TO_CREATE_TRANSACTION)

    id = result["transaction_id"]

    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = status_transaction_send(id)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0
    assert status_result["transaction_status"] == "created"

    cancel_transaction_send(id)

    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = status_transaction_send(id)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0
    assert status_result["transaction_status"] == "canceled"

def test_accept_transaction():
    result = create_transaction_send(DATA_TO_CREATE_TRANSACTION)

    id = result["transaction_id"]

    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = status_transaction_send(id)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0
    assert status_result["transaction_status"] == "created"

    accept_transaction_send(id, DATA_TO_ACCEPT_TRANSACTION)

    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = status_transaction_send(id)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0
    assert status_result["transaction_status"] == "processed"

def test_edit_retrieve_transaction():
    result = create_transaction_send(DATA_TO_CREATE_TRANSACTION)

    id = result["transaction_id"]

    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = status_transaction_send(id)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0
    assert status_result["transaction_status"] == "created"

    assert retrieve_transaction_send(id) == DATA_TO_CREATE_TRANSACTION