FROM golang:1.25.10-alpine AS base
LABEL author="Kenley Wang"


# build package
FROM base AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o clio main.go


# run in production
FROM base AS prod

WORKDIR /app

COPY --from=builder /app/clio .
COPY --from=builder /app/config/config.yaml /etc/clio/config.yaml
COPY --from=builder /app/LICENSE ./LICENSE
COPY --from=builder /app/docker/entrypoint.sh /entrypoint.sh

EXPOSE 9899
ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]