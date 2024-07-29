## API Documentation for IP Management Service

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

- **URL**: `GET /api/action-logs`
- **Description**: Retrieve logs based on various filters.
- **Query Parameters**:
  - `user_id` (optional): Filter logs by a specific user ID.
  - `before` (optional): Filter logs before a specific date in UTC+6.
  - `after` (optional): Filter logs after a specific date in UTC+6.
  - `start` (optional): Filter logs with a start date range in UTC+6.
  - `end` (optional): Filter logs with an end date range in UTC+6.

#### Example Request

```http
GET /api/action-logs?user_id=123&before=2024-01-04T00:00:00+06:00&after=2023-01-02T00:00:00+06:00&start=2023-01-02T00:00:00+06:00&end=2024-09-09T00:00:00+06:00
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
