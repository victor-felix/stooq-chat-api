# REST API Stooq Chat Application

This is an example of a REST API for chat using Websocket.

## Executing with docker-compose

    docker-compose up --build

NOTE: This way, RabbitMQ may take a while to finish execution.

## Executing main.go

    go run main.go

NOTE: This way, you have to install MongoDB and RabbitMQ, and check if config file is ready.

# REST API

The REST API to the example app is described below.

## Create an account

### Request

`POST /users`

    curl --location 'http://localhost:8080/users' \ --header 'Content-Type: application/json' \ --data-raw '{ "user_name": "Victor", "email": "victor@email.com", "password":  "Qwe12345" }'

### Response

    HTTP/1.1 201 Created
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 2

    {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}

## Send new message

### Request

`POST /messages`

    curl --location 'http://localhost:8080/messages' \ --header 'Content-Type: application/json' \ --header 'Authorization: Bearer token' \ --data-raw '{ "room_id": "507f1f77bcf86cd799439011", "user_id": "507f1f77bcf86cd799439011", "content":  "message content" }'

### Response

    HTTP/1.1 201 Created
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 201 Created
    Connection: close
    Content-Type: application/json
    Location: /thing/1
    Content-Length: 36

    { "id": "507f1f77bcf86cd799439011", "room_id": "507f1f77bcf86cd799439011", "user_id": "507f1f77bcf86cd799439011", "content":  "message content", "created_at": "2023-05-01T09:35:14.0000" }

## List rooms

### Request

`GET /rooms`

    curl --location 'http://localhost:8080/rooms' \ --header 'Content-Type: application/json' \ --header 'Authorization: Bearer token'

### Response

    HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 36

    [{"id":1,"name":"Foo"}]

## Create new room

### Request

`POST /rooms`

    curl --location 'http://localhost:8080/messages' \ --header 'Content-Type: application/json' \ --header 'Authorization: Bearer token' \ --data-raw '{ "name":  "Bar" }'

### Response

    HTTP/1.1 201 CREATED
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 201 CREATED
    Connection: close
    Content-Type: application/json
    Location: /thing/1
    Content-Length: 36

    {"id":1,"name":"Bar"}

## List room messages

### Request

`GET /room/:id/messages`

    curl --location 'http://localhost:5001/rooms/6554bb61d934ff6198d7a4e8/messages' \ --header 'Authorization: Bearer token'

### Response

    HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:31 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Location: /thing/2
    Content-Length: 35

    [{ "id": "507f1f77bcf86cd799439011", "room_id": "507f1f77bcf86cd799439011", "user_id": "507f1f77bcf86cd799439011", "content":  "message content", "created_at": "2023-05-01T09:35:14.0000" }]

## Authenticate

### Request

`POST /auth`

    curl --location 'http://localhost:5001/auth' \ --header 'Content-Type: application/ json' \ --data-raw '{ "email": "example@email.com", "password": "example"}'

### Response

    HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:31 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 74

    {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}

## Adictional Repositories

[Stooq Chat Web App](https://github.com/victor-felix/stooq-chat-webapp)

[Stooq Bot Broker](https://github.com/victor-felix/stooq-bot-broker)
