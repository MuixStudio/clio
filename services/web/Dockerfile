FROM golang:1.25.3-alpine AS base

# Install build dependencies
RUN apk add --no-cache git


# Build stage
FROM base AS builder

# 定义构建参数
ARG SERVICE_NAME=unknown-service
ARG SERVICE_VERSION=0.0.0
ARG SERVICE_DESC="No description"
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown
ARG GO_VERSION=unknown
ARG METADATA_PKG="github.com/muixstudio/clio/services/common/metadata"

WORKDIR /clio

COPY . .

RUN go mod download

# # 构建时注入元数据
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags "\
        -X '${METADATA_PKG}.Name=${SERVICE_NAME}' \
        -X '${METADATA_PKG}.Version=${SERVICE_VERSION}' \
        -X '${METADATA_PKG}.Description=${SERVICE_DESC}' \
        -X '${METADATA_PKG}.BuildTime=${BUILD_TIME}' \
        -X '${METADATA_PKG}.GitCommit=${GIT_COMMIT}' \
        -X '${METADATA_PKG}.GoVersion=${GO_VERSION}'" \
    -o app services/web/cmd/main.go


# Production stage
FROM alpine:latest AS prod
LABEL author="Kenley Wang"
LABEL maintainer="Kenley Wang"

# 添加必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /clio

# 复制二进制文件和配置
COPY --from=builder /clio/app .
COPY --from=builder /clio/services/web/etc/config.yaml /etc/config.yaml
COPY --from=builder /clio/services/web/docker/entrypoint.sh /entrypoint.sh

EXPOSE 5217

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]
