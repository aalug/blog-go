# run up migrations, user details based on docker-compose.yml
migrate_up:
	migrate -path db/migrations -database "postgresql://devuser:admin@db:5432/blog_go_db?sslmode=disable" -verbose up

# run down migrations, user details based on docker-compose.yml
migrate_down:
	migrate -path db/migrations -database "postgresql://devuser:admin@db:5432/blog_go_db?sslmode=disable" -verbose down

# generate db related go code with sqlc
sqlc:
	cmd.exe /c "docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate"
