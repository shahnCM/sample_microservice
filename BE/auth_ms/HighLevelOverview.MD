### Auth Microservice Database Schema

#### Users Table

- id: (primary key)
- username: string (unique)
- password_hash: string
- email: string (unique)
- role: enum ('standard', 'admin')

#### Sessions Table

- id: (primary key)
- started_at: datetime
- expired_at: datetime (nullable)
- refresh_count: integer
- user_id: ULID (foreign key to users table)

#### Tokens Table

- id: (primary key)
- user_id: (foreign key to users table)
- session_id: (foreign key to sessions table)
- jwt_token: string
- refresh_token: string
- issued_at: datetime
- expires_at: datetime
- token_status: enum ('fresh', 'refreshed', 'revoked')

#### Relations

- User has many Sessions (In terms of db relations)
- Session has many Tokens (In terms of db relations)

### Additional Notes

#### Token Issuing Logic

- When a user has a running session, we cannot issue a new token/start a new session.
- If a new token has to be issued/a new session needs to be started:
  - current session needs to be ended
  - current token has to be revoked

#### Token Refreshing Logic

- When refreshing a token, ensure that the token's token_status is fresh.
- A token can only be refreshed if the associated session's expired_at is null.

#### Token Revocation

- Revoked tokens are marked with token_status set to revoked.
- The associated session’s expired_at is updated to the current time.

#### User Roles

- Users have roles (standard, admin), and roles determine access control.

### Microservices Architecture

- **Auth Microservice:** Manages authentication, token issuance, token refreshing, token revocation, and session management.
- **Log Microservice:** Receives and stores user action data sent from the auth microservice through RabbitMQ.
- **IP Management Microservice:** Manages IP addresses and their associated labels/comments.
- **Gateway Microservice:** Acts as the entry point for requests, routing them to the appropriate services.

#### Communication

- The auth microservice will send user action data to the log microservice through RabbitMQ.