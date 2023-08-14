# Go Blog 

### Built in Go 1.20

### The app uses:
- Postgres
- Docker
- [Gin](https://github.com/gin-gonic/gin)
- gRPC
- [asynq](https://github.com/hibiken/asynq)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://github.com/kyleconroy/sqlc)
- [testify](https://github.com/stretchr/testify)
- [PASETO Security Tokens](https://github.com/o1egl/paseto)

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
 - `/users` - handles POST requests to create users.
 - `/users/login` - handles POST requests to log in users.

### Tokens/Session
 - `/tokens/renew` - handles  POST requests to renew the access tokens.

#### Category
 - `/category` - handles POST requests to create categories
 - `/category/{name}` - handles DELETE requests to delete a category

### Posts
- `/posts` - handles POST requests to create posts.
- `/posts/{id}` - handles DELETE requests to delete a post.
- `/posts/id/{id}` and `/posts/title/{slug}` - handles GET requests to get post details.
- `/posts/all` - handles GET requests to list all posts. Query params: `page`, `page_size`.
- `/posts/author` - handles GET requests to list posts created by author 
with given name (that username or email contain given string). 
Query params: `page`, `page_size`, `author`.
- `/posts/category` - handles GET requests to list posts from the given category.
Query params: `page`, `page_size`, `category_id`.
- `/posts/tags` - handles GET requests to list posts with given tags.
Query params: `page`, `page_size`, `tag_ids` where `tag_ids` is 
comma-separated int format (e.g. `&tag_ids=1,2,3`).
- `/posts/{id}` - handles PATCH requests to update the post.

### Comments
- `/comments` - handles POST requests to create a comment.
- `/comments/{id}` - handles DELETE requests to delete a comment.
- `/comments/{id}` - handles PATCH requests to update a comment.
- `/comments/{post_id}` - handles GET requests to list comments of a post.
Query params: `page` and `page_size`.

## Documentation
### API
The API (HTTP gateway) documentation can be found at
[this swaggerhub page](https://app.swaggerhub.com/apis/AAGULCZYNSKI/blog-go/1.0)
and (after running the server) at http://localhost:8080/docs/

### Database
The database's schema and intricate details can be found on 
dedicated webpage, which provides a comprehensive overview 
of the data structure, tables, relationships, and other essential 
information. To explore the database further, please visit
this [dbdocs.io webpage](https://dbdocs.io/aalug/blog_go).
Password: `bloggopassword`
