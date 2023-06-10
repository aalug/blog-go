FROM golang:1.20-alpine3.18 AS builder

WORKDIR /app
COPY . .
RUN go build -o main main.go

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.1/migrate.linux-amd64.tar.gz  | tar xvz

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

COPY wait-for.sh .
COPY start.sh .

#COPY dn/migrations ./migrations

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT ["/app/start.sh"]