## Installation

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