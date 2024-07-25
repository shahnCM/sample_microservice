## Authentication Service API Documentation

### 1. Register a New User

**Endpoint**: `POST /v1/auth/token/register`

**Description**: Registers a new user in the system.

**Request**:
```bash
curl --location 'localhost:8060/v1/auth/token/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "mahir",
    "password": "abc123456",
    "password_confirm": "abc123456",
    "email": "not@null.com"
}'
```

**Request Body**:
- `username` (string): The username of the new user.
- `password` (string): The password of the new user.
- `password_confirm` (string): The password confirmation.
- `email` (string): The email address of the new user.

### 2. Authenticate and Get Fresh Token

**Endpoint**: `POST /v1/auth/token/fresh`

**Description**: Authenticates a user and returns a fresh JWT token.

**Request**:
```bash
curl --location 'localhost:8060/v1/auth/token/fresh' \
--header 'Content-Type: application/json' \
--data '{
    "username": "mahir",
    "password": "abc123456"
}'
```

**Request Body**:
- `username` (string): The username of the user.
- `password` (string): The password of the user.

### 3. Verify Token

**Endpoint**: `PUT /v1/auth/token/verify`

**Description**: Verifies the provided JWT token.

**Request**:
```bash
curl --location --request PUT 'localhost:8060/v1/auth/token/verify' \
--header 'Authorization: Bearer <jwt_token>' \
--data ''
```

**Headers**:
- `Authorization`: Bearer token string.

### 4. Revoke Token

**Endpoint**: `PUT /v1/auth/token/revoke`

**Description**: Revokes the provided JWT token.

**Request**:
```bash
curl --location --request PUT 'localhost:8060/v1/auth/token/revoke' \
--header 'Authorization: Bearer <jwt_token>' \
--data ''
```

**Headers**:
- `Authorization`: Bearer token string.

### 5. Refresh Token

**Endpoint**: `POST /v1/auth/token/refresh`

**Description**: Refreshes the JWT token using a refresh token.

**Request**:
```bash
curl --location 'localhost:8060/v1/auth/token/refresh' \
--header 'Authorization: Bearer <jwt_token>' \
--header 'Content-Type: application/json' \
--data '{
    "type": "refresh",
    "token": "<refresh_token>"
}'
```

**Headers**:
- `Authorization`: Bearer token string.
- `Content-Type`: application/json

**Request Body**:
- `type` (string): Should be "refresh".
- `token` (string): The refresh token.

### 6. Get User Sessions

**Endpoint**: `GET /v1/sessions/users/{userId}`

**Description**: Retrieves the sessions of a specific user.

**Request**:
```bash
curl --location 'localhost:8060/v1/sessions/users/1?page=1&per_page=3' \
--header 'Authorization: Bearer <jwt_token>'
```

**Headers**:
- `Authorization`: Bearer token string.

**Query Parameters**:
- `page` (int): The page number.
- `per_page` (int): The number of sessions per page.

---

Replace `<jwt_token>` and `<refresh_token>` with the actual tokens you have. This documentation provides details on how to use each endpoint, including the required headers and request body data.