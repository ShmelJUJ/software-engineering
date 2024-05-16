import json
from urllib.request import urlopen, Request

TRANSACTION_URL = "http://localhost:8083"
API_PATH = "/api/v1"
CREATE_TRANSACTION_PATH = "/transaction/create"
LOGIN_PATH = "/transaction/login"


CREATOR_DATA = {
    "email": "luffy@example.com",
    "password": "pirateking"
}

BUYER_DATA = {
    "email": "nami@example.com",
    "password": "treasure"
}

DATA_TO_CREATE_TRANSACTION = {
    "money_info": {
        "amount": 1,
        "currency": "string",
        "method": "algorand"
    },
    "receiver": {
        "user_id": "85e6a060-f914-48d1-b73a-23b7e6c81f46",
        "wallet_id": "3735b92d-5dcb-4dc0-a8b0-54415d0c52d3"
    }
}

DATA_TO_ACCEPT_TRANSACTION = {
    "sender": {
        "user_id": "65a8ed73-b6f3-4543-82a6-7ab9ef6e9c7b",
        "wallet_id": "c09e7795-8b25-4639-b7e3-8c2592298eb8"
    }
}

DATA_TO_RETRIEVE_TRANSACTION = {
    "money_info": {
        "amount": 2,
        "currency": "string",
        "method": "algorand"
    }
}

EDITED_TRANSACTION_TRANSACTION = {
    "amount": 2,
    "currency": "string",
    "method": "algorand",
    "receiver": {
        "user_id": "85e6a060-f914-48d1-b73a-23b7e6c81f46",
        "wallet_id": "3735b92d-5dcb-4dc0-a8b0-54415d0c52d3"
    },
    "status": "created"
}

ACCEPT_DATA = {
    "sender": {
        "user_id": "65a8ed73-b6f3-4543-82a6-7ab9ef6e9c7b",
        "wallet_id": "c09e7795-8b25-4639-b7e3-8c2592298eb8"
    }
}

def headers_auth(token):
    return {'content-type': 'application/json',
            'Authorization': token
            }

def headers():
    return {'content-type': 'application/json'}

def create_transaction_send(data: dict, token: str) -> dict:
    req = Request(TRANSACTION_URL + API_PATH + CREATE_TRANSACTION_PATH, data=json.dumps(
        data).encode(), headers=headers_auth(token), method='POST')
    response = urlopen(req)
    assert response.getcode() == 200
    result = json.loads(response.read().decode())
    return result

def get_auth_token_send(data: list) -> dict:
    req = Request(TRANSACTION_URL + API_PATH + LOGIN_PATH, data=json.dumps(
        data).encode(), headers=headers(), method='POST')
    response = urlopen(req)
    assert response.getcode() == 200
    result = json.loads(response.read().decode())
    return result["auth_token"]

def status_transaction_send(id: str, token: str) -> list:
    req = Request(TRANSACTION_URL + API_PATH + "/transaction/" + id + "/retrieve/status", headers=headers_auth(token),
                  data=json.dumps(dict({"transaction_id": id})).encode(), method='GET')
    response = urlopen(req)
    assert response.getcode() == 200
    return json.loads(response.read().decode())

def cancel_transaction_send(id: str, token: str):
    req = Request(TRANSACTION_URL + API_PATH + "/transaction/" + id + "/cancel", headers=headers_auth(token),
                  data=json.dumps(dict({"transaction_id": id})).encode(), method='POST')
    response = urlopen(req)
    assert response.getcode() == 200

def accept_transaction_send(id: str, data: dict, token: str):
    req = Request(TRANSACTION_URL + API_PATH + "/transaction/" + id + "/accept", headers=headers_auth(token),
                  data=json.dumps(data).encode(), method='POST')
    response = urlopen(req)
    assert response.getcode() == 200

def retrieve_transaction_send(id: str, token: str):
    req = Request(TRANSACTION_URL + API_PATH + "/transaction/" + id + "/retrieve", headers=headers_auth(token), method='GET')
    response = urlopen(req)
    assert response.getcode() == 200
    return json.loads(response.read().decode())

def edit_transaction(data: dict, id: str, token: str):
    req = Request(TRANSACTION_URL + API_PATH + "/transaction/" + id + "/edit", data=json.dumps(
        data).encode(), headers=headers_auth(token), method='POST')
    response = urlopen(req)
    assert response.getcode() == 200