### Auth Microservice High-Level Overview

#### **Database Schema**

- **Users Table**
  - **id**: (primary key)
  - **username**: string (unique)
  - **password_hash**: string
  - **email**: string (unique)
  - **role**: enum ('standard', 'admin')
  - **SessionTokenTraceId**: *string (unique, nullable)
  - **LastSessionId**: *uint (unique, nullable)
  - **LastSession**: *Session (foreign key)

- **Sessions Table**
  - **id**: (primary key)
  - **user_id**: ULID (foreign key to users table)
  - **session_token_trace_id**: *string (not null)
  - **refresh_count**: integer
  - **starts_at**: datetime (not null)
  - **ends_at**: datetime (nullable)
  - **revoked**: boolean
  - **refresh_ends_at**: datetime
  - **created_at**: datetime
  - **updated_at**: datetime


#### **Relations**

- A **User** can have many **Sessions**.
- A **Session** can have many **Tokens** (though tokens are not used directly in actions).

#### **Action Logics**

- **Register**: Hashes the password and stores the user information in the database.

- **Login**:
  - Verifies username and password.
  - Creates a new session and updates the user with the new session ID and token ID.
  - Ends any previous active sessions.

- **Revoke**:
  - Verifies the JWT token.
  - Updates the associated session's `expires_at` to the current time.
  - Marks the token's status as 'revoked'.

- **Verify**:
  - Verifies the JWT token and returns user information if valid.

- **Refresh**:
  - Verifies both JWT and refresh tokens concurrently.
  - Ensures JWT token is expired and refresh token is valid.
  - Checks consistency between JWT and refresh tokens.
  - Locks user and session records for update.
  - Generates new JWT and refresh tokens.
  - Updates session with new token details and expiration times.
  - Updates user with the new session token trace ID.
  - Rolls back transaction on error and commits if successful.

#### **Microservices Architecture**

- **Auth Microservice**: Handles authentication, token issuance, token refreshing, token revocation, and session management.
- **IP Management Microservice**: Manages IP addresses and their associated labels/comments.
- **Gateway Microservice**: Acts as the entry point for requests, routing them to the appropriate services.

#### **Communication**

- The auth microservice sends user action data to the log microservice through RabbitMQ.

#### **Endpoints**

- `/fresh`: Issues a new token.
- `/refresh`: Refreshes an expired token.
- `/revoke`: Revokes a token.
- `/verify`: Verifies a token.
- `/register`: Registers a new user.
- `/sessions/users/:id`: Retrieves session details for a specific user (admin only).


