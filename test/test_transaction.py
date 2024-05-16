from time import sleep
import requsts as r


def test_creation_transaction():
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

def test_edit_retrieve_transaction():
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

    r.edit_transaction(r.DATA_TO_RETRIEVE_TRANSACTION, id, token)

    assert r.retrieve_transaction_send(id, token) == r.EDITED_TRANSACTION_TRANSACTION