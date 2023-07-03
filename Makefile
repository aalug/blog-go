# generate migrations, $(name) - name of the migration
generate_migrations:
	migrate create -ext sql -dir db/migrations -seq $(name)

# run up migrations, user details based on docker-compose.yml
migrate_up:
	migrate -path db/migrations -database "postgresql://devuser:admin@localhost:5432/blog_go_db?sslmode=disable" -verbose up

# run down migrations, user details based on docker-compose.yml
migrate_down:
	migrate -path db/migrations -database "postgresql://devuser:admin@localhost:5432/blog_go_db?sslmode=disable" -verbose down

# generate database documentation on the dbdocs website
db_docs:
	dbdocs build docs/database.dbml

# generate .sql file with database schema
db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/database.dbml

# generate db related go code with sqlc
sqlc:
	cmd.exe /c "docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate"

# run all tests
test:
	go test -v -cover ./...

# run tests in the given path (p) and display results in the html file
test_coverage:
	go test $(p) -coverprofile=coverage.out && go tool cover -html=coverage.out

# generate mock db for testing
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/aalug/blog-go/db/sqlc Store

# remove old files and generate new protobuf files
protoc:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=protobuf --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=blog_go \
	protobuf/*.proto

# starts just the db container
start_db:
	docker-compose logs -f db

.PHONY: generate_migrations, migrate_up, migrate_down, sqlc, test, test_coverage, mock, db_schema, db_docs, protobuf, start_db