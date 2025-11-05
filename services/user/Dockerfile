FROM golang:1.24.7-alpine AS base
LABEL author="Kenley Wang"


# build package
FROM base AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-rpc ./user.go


# run in production
FROM base AS prod

WORKDIR /app

COPY --from=builder /app/user-rpc .
COPY --from=builder /app/etc/config.yaml /etc/config.yaml
COPY docker/entrypoint.sh /entrypoint.sh

EXPOSE 5000
ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]