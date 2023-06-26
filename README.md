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
3. Rename `app.env.sample` to `app.env` and replace the values
4. Run in your terminal:
    - `docker-compose up --build` to run the containers
5. Now everything should be ready and server running on `SERVER_ADDRESS` specified in `app.env`

## Testing
1. Run the postgres container (`docker-compose up`)
2. Run in your terminal:
    - `make test` to run all tests
   or
    - `make test_coverage p={PATH}` - to get the coverage in the HTML format - where `{PATH}` is the path to the target directory for which you want to generate test coverage. The `{PATH}` should be replaced with the actual path you want to use. For example `./api`
   or
    - use standard `go test` commands (e.g. `go test -v ./api`)

## API endpoints
#### Users
 - `/users` - handles POST requests to create users
 - `/users/login` - handles POST requests to log in users

#### Category
 - `/category` - handles POST requests to create categories
 - `/category/{name}` - handles DELETE requests to delete a category

### Posts
- `/posts` - handles POST requests to create posts
- `/posts/{id}` - handles DELETE requests to delete a post
- `/posts/id/{id}` and `/posts/title/{slug}` - handles GET requests to get post details
- `/posts/all` - handles GET requests to list all posts. Query params: `page`, `page_size`
- `/posts/author` - handles GET requests to list all posts created by author 
with given name (that username or email contain given string). 
Query params: `page`, `page_size`, `author`

