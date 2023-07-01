FROM golang:1.20-alpine3.18 AS builder

WORKDIR /app
COPY . .
RUN go build -o main_grpc main_grpc.go

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.1/migrate.linux-amd64.tar.gz  | tar xvz

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main_grpc .
COPY --from=builder /app/migrate .

COPY app.env .
COPY wait-for.sh .
COPY start.sh .

COPY db/migrations ./migrations

EXPOSE 9090
CMD [ "/app/main_grpc" ]
ENTRYPOINT ["/app/start.sh"]