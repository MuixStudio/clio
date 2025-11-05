FROM golang:1.25.3-alpine AS base
LABEL author="Kenley Wang"

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
ARG GIT_BRANCH=unknown
ARG GO_VERSION=unknown

WORKDIR /clio

COPY . .

RUN go mod download

# # 构建时注入元数据
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags "\
        -X 'github.com/muixstudio/clio/services/common/metadata.Name=${SERVICE_NAME}' \
        -X 'github.com/muixstudio/clio/services/common/metadata.Version=${SERVICE_VERSION}' \
        -X 'github.com/muixstudio/clio/services/common/metadata.Description=${SERVICE_DESC}' \
        -X 'github.com/muixstudio/clio/services/common/metadata.BuildTime=${BUILD_TIME}' \
        -X 'github.com/muixstudio/clio/services/common/metadata.GitCommit=${GIT_COMMIT}' \
        -X 'github.com/muixstudio/clio/services/common/metadata.GitBranch=${GIT_BRANCH}' \
        -X 'github.com/muixstudio/clio/services/common/metadata.GoVersion=${GO_VERSION}'" \
    -o app services/web/cmd/cmd.go


# Production stage
FROM alpine:latest AS prod

# 添加必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /clio

# 复制二进制文件和配置
COPY --from=builder /clio/app .
COPY --from=builder /clio/services/web/etc/config.yaml /etc/config.yaml
COPY --from=builder /clio/services/web/docker/entrypoint.sh /entrypoint.sh

EXPOSE 3917

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]
