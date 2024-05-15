from time import sleep
import requsts as r

def test_cancel_transaction():
    token = r.get_auth_token_send(r.CREATOR_DATA)

    sleep(1)

    result = r.create_transaction_send(r.DATA_TO_CREATE_TRANSACTION, token)
    id = result["transaction_id"]

    sleep(1)
    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = r.status_transaction_send(id, token)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0
    sleep(1)
    assert status_result["transaction_status"] == "created"

    r.cancel_transaction_send(id, token)

    sleep(1)
    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = r.status_transaction_send(id, token)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0
    sleep(1)
    assert status_result["transaction_status"] == "canceled"

def test_accept_transaction():
    creator_token = r.get_auth_token_send(r.CREATOR_DATA)

    sleep(1)

    result = r.create_transaction_send(r.DATA_TO_CREATE_TRANSACTION, creator_token)

    id = result["transaction_id"]

    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = r.status_transaction_send(id, creator_token)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0
    assert status_result["transaction_status"] == "created"

    sender_token = r.get_auth_token_send(r.CREATOR_DATA)

    r.accept_transaction_send(id, r.DATA_TO_ACCEPT_TRANSACTION, sender_token)

    sleep(5)

    retries = 10
    while retries > 0:
        sleep(0.5)
        status_result = r.status_transaction_send(id, creator_token)
        if status_result != []:
            break
        else:
            retries -= 1
    assert retries > 0

    assert status_result["transaction_status"] == "succeeded"
