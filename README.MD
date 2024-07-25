# Api Docs

### Authentication Service

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