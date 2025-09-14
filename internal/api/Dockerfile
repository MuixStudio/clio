FROM golang:1.24.7-alpine AS base
LABEL author="Kenley Wang"


# build package
FROM base AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./main.go


# run in production
FROM base AS prod

WORKDIR /app

COPY --from=builder /app/api .
COPY docker/entrypoint.sh /entrypoint.sh

EXPOSE 5000
ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]