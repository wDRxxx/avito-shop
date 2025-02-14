FROM golang:1.23-alpine AS builder

COPY . /avito-shop/source/
WORKDIR /avito-shop/source/

RUN go build -o ./bin/avito-shop cmd/api/main.go
RUN go build -o ./bin/migrator cmd/migrator/main.go

FROM alpine:3.13

WORKDIR /root/
COPY --from=builder /avito-shop/source/bin/ .
COPY --from=builder /avito-shop/source/migrations /migrations/
COPY --from=builder /avito-shop/source/docker.env .

CMD ["sh", "-c", "./migrator --env-path=docker.env --migrations-path=/migrations/ && ./avito-shop --env-path=docker.env --env-level=prod --logs-path=/avito-shop/logs"]