# Go Blog 

### Built in Go 1.20

### The app uses:
- Postgres
- Docker
- [Gin](https://github.com/gin-gonic/gin)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://github.com/kyleconroy/sqlc)
- [testify](https://github.com/stretchr/testify)
- [PASETO Security Tokens](github.com/o1egl/paseto)

## Getting started
1. Clone the repository
2. Go to the project's root directory
3. Run in your terminal:
    - `docker-compose up` to run the containers
4. Now everything should be ready and API available http 

## Testing
1. Run the postgres container (`docker-compose up`)
2. Run in your terminal:
    - `make test` to run all tests
   or
    - `make test_coverage` to run all tests and see the coverage in the html format
   or
    - use stanard `go test` commands (e.g. `go test -v ./api`)

## API endpoints
#### Users
 - `/users` - handles POST requests to create users
 - `/users/login` - handles POST requests to log in users

#### Category
 - `/category` - handles POST requests to create categories