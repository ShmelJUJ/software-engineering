# **\[Core Backend API\]** 

## **Problem** 

Нужна единая точка входа в сервисы системы оплаты по qr кодам: обновление статуса платежа, взаимодействие клиентов с другими сервисами

**Anti-Goals**

Генерация и сканирование кр кодов, работа с api банков, хранение пд

## **Solution**

Создание сервиса, с которым будут взаимодействовать клиентыи который будет хранить информацию о платежах 

## **Changes and Additions to Public Interfaces**

## **API**
**\[Payment\]**

POST ```api/v1/payment/create```

headers

x_auth_token

request_body
```
{
    "idempotency_token": "uuid4",
    "money_info": {
        "method_id": "recipcient_payment_method_id",
        "currency": "BTC"
        "amount": 1.0
    },
    "recipcient": {
        "recipcient_id": "string"
    },
    "sender": {
        "sender_id": "string"
    }
}
```

response_body
200

```
{
    "transaction_id": "sfdsagfs"
}
```
400

```
{
    "code": "validation_error", //machine_redable_code
    "message": "validation error reason"
}
```
403

{
    "code": "acces_error", //machine_redable_code
    "message": "aboba"
}

500
...

GET ```api/v1/payment/{payment_id}/retrieve```

headers

x_auth_token

response 200
```
{
    payment: {
        sender:{
            user_id: string,
            method_id: string
        },
        recipcient: {
            user_id: string,
            method_id: string
        },
        currency: ISOFORMAT
        amount: double,
        status: enum
    }
}
```
404
{
    code: "not_found"
    message: "payment {payment_id} not found"
}
400, 403, 500

GET ```api/v1/payment/{payment_id}/retrieve/status```

headers

x_auth_token

response 200
```
{
    payment_status: enum
}
```
404
{
    code: "not_found"
    message: "payment {payment_id} not found"
}
400, 403, 500

POST ```api/v1/payment/{id}/cancel```

{
    "idempotency_token": uuid
    "reason" enum
}

200 - ok
400 - validation error or status is bank_processing
500

POST ```api/v1/payment/{id}/edit```

{
    "idempotency_token": uuid,
    "money_info": {
        "method_id": "recipcient_payment_method_id",
        "currency": "BTC"
        "amount": 1.0
    },
}

200 - ok
400 - validation error or status is bank_processing
500

POST POST ```api/v1/payment/{id}/accept```

{
    "idempotency_token": uuid,
    "accept_info":{
        // TODO узгать у Соболева что тут будет лежать
    }
}

200 - ok
400 - validation error
500