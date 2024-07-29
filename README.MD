# Installation

### Authentication Service

- Clone the repo
- Go to dir: `BE/auth_ms/codes`
- Create a database for this project in your environment
- Rename example.env to .env
- Configure environment properly
- Run `go mod tidy` in cli
- Run `RUN_MIGRATION=TRUE go run main.go` if you want to run migration
- Otherwise only run `go run main.go`
- This app will run on `localhost:<port>`

### IP Management Service

- Clone the repo
- Go to dir: `BE/ip_management_ms/codes`
- Create a database for this project in your environment
- Rename example.env to .env
- Configure environment properly
- Run `composer install` in cli
- Run `php artisan key:generate` in cli
- Run `php artisan migrate` in cli
- Run `php artisan serve` in cli
- This app will run on `localhost:<port>`

# Api Docs

## Authentication Service

#### Base URL

The base URL for all API requests is:

`http://localhost:8060/auth`


#### Health Check

- **GET** `/health-check`

**Description:** Checks if the server is alive.

**Response:**
    - Status: `200 OK`
    - Body: `{"status":"Server is alive"}`

---
#### Authentication and Token Management

##### Register User


- **POST** `/api/v1/token/register`

**Description:** Registers a new user.

**Request:**

```sh
curl --location 'http://localhost:8060/auth/api/v1/token/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "mahir",
    "password": "abc123456",
    "password_confirm": "abc123456",
    "email": "not@null.com"
}'
```

**Response:**

```json
{
    "status": "Success",
    "status_code": 201,
    "data": null,
    "meta": null,
    "links": null
}
```

##### Fresh Token


- **POST** `/api/v1/token/fresh`

**Description:** Issues a fresh JWT and refresh token.

**Request:**

```sh
curl --location 'http://localhost:8060/auth/api/v1/token/fresh' \
--header 'Content-Type: application/json' \
--data '{
    "username": "mahir",
    "password": "abc123456"
}'
```

**Response:**

```json
{
    "status": "Success",
    "status_code": 200,
    "data": {
        "jwt": {
            "type": "JWT",
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJ1c2VyX3JvbGUiOiJhZG1pbiIsInRva2VuX2lkIjoiJDJhJDEwJC80TWIuZzVrdElXNjNhV3ZXLzlRcU9tVjFickVSMEdmNHZMYjJISGRoV0pqMU9CN0txN1dxIiwiZXhwIjoxNzIxODM4NjQzfQ.WrIvWwi396VsI9lVJDy14rAafU5j8ikjOcikSv-g49Y",
            "expires_at": 1721838643
        },
        "refresh": {
            "type": "REFRESH",
            "token": "eyJ1c2VyX2lkIjoyLCJ1c2VyX3JvbGUiOiJhZG1pbiIsInRva2VuX2lkIjoiJDJhJDEwJC80TWIuZzVrdElXNjNhV3ZXLzlRcU9tVjFickVSMEdmNHZMYjJISGRoV0pqMU9CN0txN1dxIiwiZXhwIjoxNzIyNDM5ODQzfQ.enmYQKamekbfSB57FTcC9v-MMa1fI-8Jt40cGhU0Ekk",
            "expires_at": 1722439843
        },
        "user": {
            "role": "admin",
            "username": "mahirA"
        }
    },
    "meta": null,
    "links": null
}
```

##### Verify Token


- **PUT** `/api/v1/token/verify`

**Description:** Verifies the provided JWT.

**Request:**

```sh
curl --location --request PUT 'http://localhost:8060/auth/api/v1/token/verify' \
--header 'Authorization: Bearer <jwt_token>' \
--data ''
```

**Response:**

```json
{
    "status": "Success",
    "status_code": 200,
    "data": {
        "role": "standard",
        "username": "mahir"
    },
    "meta": null,
    "links": null
}
```

##### Revoke Token


- **PUT** `/api/v1/token/revoke`

**Description:** Revokes the provided JWT.

**Request:**

```sh
curl --location --request PUT 'http://localhost:8060/auth/api/v1/token/revoke' \
--header 'Authorization: Bearer <jwt_token>' \
--data ''
```

**Response:**
- Status: `204 No Content`

##### Refresh Token


- **POST** `/api/v1/token/refresh`

**Description:** Refreshes the provided refresh token and issues new JWT and refresh token.

**Request:**

```sh
curl --location 'http://localhost:8060/auth/api/v1/token/refresh' \
--header 'Authorization: Bearer <jwt_token>' \
--header 'Content-Type: application/json' \
--data '{
    "type": "refresh",
    "token": "<refresh_token>"
}'
```

**Response:**

```json
{
    "status": "Success",
    "status_code": 200,
    "data": {
        "jwt": {
            "type": "JWT",
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJ1c2VyX3JvbGUiOiJhZG1pbiIsInRva2VuX2lkIjoiJDJhJDEwJExVdml1dUF5dTEudS5UYUxmckkxcU90d01zS2hqS3o3d1FtM0g2QVJIeENycktDaER6ZW1hIiwiZXhwIjoxNzIxOTA0OTg4fQ.4qaaHWsi0S_2edUIpEaiEzrE5Hb16bBEmFYeuyfLlOk",
            "expires_at": 1721904988
        },
        "refresh": {
            "type": "REFRESH",
            "token": "eyJ1c2VyX2lkIjoyLCJ1c2VyX3JvbGUiOiJhZG1pbiIsInRva2VuX2lkIjoiJDJhJDEwJExVdml1dUF5dTEudS5UYUxmckkxcU90d01zS2hqS3o3d1FtM0g2QVJIeENycktDaER6ZW1hIiwiZXhwIjoxNzIyNTA2MTg4fQ.YfKfcaohhXGnnJ_BhBDyFATcofTsqAFQecAcPQ3ET8Y",
            "expires_at": 1722506188
        }
    },
    "meta": null,
    "links": null
}
```

---
#### Session Data

##### User Sessions (Admin Only)


- **GET** `/api/v1/sessions/users/{user_id}?page=1&per_page=3`

**Description:** Retrieves user sessions (admin only).

**Request:**

```sh
curl --location 'http://localhost:8060/auth/api/v1/sessions/users/1?page=1&per_page=3' \
--header 'Authorization: Bearer <jwt_token>'
```

**Response:**

```json
{
    "status": "Success",
    "status_code": 200,
    "data": [
        {
            "Id": 7,
            "UserId": 1,
            "SessionTokenTraceId": "01J33SYE87ET9FQ5HB3DKA3KST",
            "RefreshCount": 0,
            "StartsAt": "2024-07-19T02:53:39.26+06:00",
            "EndsAt": "2024-07-19T03:00:48.555+06:00",
            "Revoked": true,
            "RefreshEndsAt": "2024-07-19T03:00:48.555+06:00",
            "CreatedAt": "2024-07-19T02:53:39.26+06:00",
            "UpdatedAt": "2024-07-19T03:00:48.556+06:00"
        },
        {
            "Id": 8,
            "UserId": 1,
            "SessionTokenTraceId": "01J33TBHFCRSESRKS57889KXQ2",
            "RefreshCount": 0,
            "StartsAt": "2024-07-19T03:00:48.542+06:00",
            "EndsAt": "2024-07-21T19:55:54.583+06:00",
            "Revoked": true,
            "RefreshEndsAt": "2024-07-21T19:55:54.583+06:00",
            "CreatedAt": "2024-07-19T03:00:48.543+06:00",
            "UpdatedAt": "2024-07-21T19:55:54.583+06:00"
        }
    ],
    "meta": {
        "current_page": 3,
        "from": 7,
        "last_page": 34,
        "limit": 3,
        "path": "http://localhost:8060/v1/sessions/users/1",
        "to": 9,
        "total": 100
    },
    "links": {
        "first": "http://localhost:8060/v1/sessions/users/1&page=1&per_page=3",
        "last": "http://localhost:8060/v1/sessions/users/1&page=34&per_page=3",
        "next": "http://localhost:8060/v1/sessions/users/1&page=4&per_page=3",
        "prev": "http://localhost:8060/v1/sessions/users/1&page=2&per_page=3"
    }
}
```

## IP Management Service

### Get All IP Addresses

### Base URL

The base URL for all API requests is:

`http://localhost:8000/ip_management`

**Endpoint:**
```
GET /api/v1/ip_addresses/
```

**Headers:**
```http
Authorization: Bearer <JWT>
```

**Response:**
```json
{
    "status": "Success",
    "status_code": 200,
    "message": "Ip Address",
    "data": [
        {
            "ip_address": "127.0.0.10",
            "label": "Home Updated",
            "created_at": "2024-07-29T17:00:27.000+06:00"
        },
        {
            "ip_address": "127.0.0.2",
            "label": "Office",
            "created_at": "2024-07-29T21:43:17.000+06:00"
        },
        {
            "ip_address": "127.0.0.4",
            "label": "Office Mohakhali Updated 4",
            "created_at": "2024-07-29T21:46:03.000+06:00"
        },
        {
            "ip_address": "127.0.0.5",
            "label": "Office Mohakhali 4",
            "created_at": "2024-07-29T21:51:49.000+06:00"
        },
        {
            "ip_address": "127.0.0.6",
            "label": "Office Mohakhali 6",
            "created_at": "2024-07-29T21:52:22.000+06:00"
        }
    ],
    "links": {
        "first": "http://localhost:8000/ip_management/api/v1/ip_addresses?page=1",
        "last": "http://localhost:8000/ip_management/api/v1/ip_addresses?page=1",
        "prev": null,
        "next": null
    },
    "meta": {
        "current_page": 1,
        "from": 1,
        "last_page": 1,
        "links": [
            {
                "url": null,
                "label": "&laquo; Previous",
                "active": false
            },
            {
                "url": "http://localhost:8000/ip_management/api/v1/ip_addresses?page=1",
                "label": "1",
                "active": true
            },
            {
                "url": null,
                "label": "Next &raquo;",
                "active": false
            }
        ],
        "path": "http://localhost:8000/ip_management/api/v1/ip_addresses",
        "per_page": 15,
        "to": 5,
        "total": 5
    }
}
```

### Get Specific IP Address

**Endpoint:**
```
GET /api/v1/ip_addresses/{id}
```

**Headers:**
```http
Authorization: Bearer <JWT>
```

**Response:**
```json
{
    "status": "Success",
    "status_code": 200,
    "message": "",
    "data": {
        "ip_address": "127.0.0.4",
        "label": "Office Mohakhali Updated 4",
        "created_at": "2024-07-29T21:46:03.000+06:00"
    }
}
```

### Action Log Search

#### Endpoint

- **URL**: `GET /api/v1/action-logs`
- **Description**: Retrieve logs based on various filters.
- **Query Parameters**:
  - `user_id` (optional): Filter logs by a specific user ID.
  - `before` (optional): Filter logs before a specific date in UTC+6.
  - `after` (optional): Filter logs after a specific date in UTC+6.
  - `start` (optional): Filter logs with a start date range in UTC+6.
  - `end` (optional): Filter logs with an end date range in UTC+6.

#### Example Request

```http
GET /api/v1/action-logs?user_id=123&before=2024-01-04T00:00:00+06:00&after=2023-01-02T00:00:00+06:00&start=2023-01-02T00:00:00+06:00&end=2024-09-09T00:00:00+06:00
```

#### Example Response

```json
[
  {
    "id": 1,
    "user_id": 123,
    "username": "exampleUser",
    "action": "store",
    "change": {
      "previous_data": {},
      "new_data": {
        "ip": "192.168.1.1",
        "label": "Office Network"
      }
    },
    "logged_at": "2024-01-02T00:00:00+06:00"
  },
  {
    "id": 2,
    "user_id": 123,
    "username": "exampleUser",
    "action": "update",
    "change": {
      "previous_data": {
        "ip": "192.168.1.1",
        "label": "Office Network"
      },
      "new_data": {
        "ip": "192.168.1.2",
        "label": "Office Network - Updated"
      }
    },
    "logged_at": "2024-01-03T00:00:00+06:00"
  }
]
```

### Create New IP Address

**Endpoint:**
```
POST /api/v1/ip_addresses/
```

**Headers:**
```http
Authorization: Bearer <JWT>
Content-Type: application/json
```

**Request Body:**
```json
{
    "ip": "127.0.0.6",
    "label": "Office Mohakhali 6"
}
```

**Success Response:**
```json
{
    "status": "Success",
    "status_code": 201,
    "message": "Ip Address Create Successful",
    "data": {
        "ip_address": "127.0.0.9",
        "label": "Office Mohakhali 6",
        "created_at": "2024-07-30T02:18:23.000+06:00"
    }
}
```

**Error Response:**
```json
{
    "status": "Error",
    "status_code": 422,
    "message": "Validation Error",
    "validation_errors": {
        "ip": [
            "The ip field must be a valid IP address."
        ]
    }
}
```

### Update IP Address

**Endpoint**: `PUT /api/v1/ip_addresses/{id}`

**Description**: Update an existing IP address.

**Path Parameters**:
- `id`: The ID of the IP address to update.

**Request Headers**:
- `Authorization`: Bearer token for authorization.

**Request Body**:
- `ip` (optional): The new IP address.
- `label` (optional): The new label for the IP address.

**Example Request**:

```bash
curl --location --request PUT 'http://localhost:8000/ip_management/api/v1/ip_addresses/6' \
--header 'Authorization: Bearer <your_jwt_token>' \
--header 'Content-Type: application/json' \
--data '{
  "ip": "127.0.0.23",
  "label": "Office Mohakhali Updated 4"
}'
```

**Example Response**:

```json
{
  "status": "Success",
  "status_code": 201,
  "message": "Ip Address Update Successful",
  "data": null
}
```

**Response Codes**:
- `201 Created`: The IP address was successfully updated.
- `400 Bad Request`: The request was invalid.
- `401 Unauthorized`: Authorization token is missing or invalid.
- `404 Not Found`: The specified IP address ID does not exist.
- `500 Internal Server Error`: An error occurred on the server.
