DOCKER_REGISTRY=muixstudio
WEB_IMAGE=$(DOCKER_REGISTRY)/clio-web
API_IMAGE=$(DOCKER_REGISTRY)/clio-api
USER_RPC_IMAGE=$(DOCKER_REGISTRY)/clio-user-rpc
VERSION=0.0.1

build-web:
	@echo "Building web image..."
	docker buildx build -t $(WEB_IMAGE):$(VERSION) --platform=linux/amd64,linux/arm64 ./web
	@echo "$(WEB_IMAGE):$(VERSION) built successfully."

build-api:
	@echo "Building web image..."
	docker buildx build -t $(API_IMAGE):$(VERSION) --platform=linux/amd64,linux/arm64 ./server/api
	@echo "$(API_IMAGE):$(VERSION) built successfully."

build-user-rpc:
	@echo "Building web image..."
	docker buildx build -t $(USER_RPC_IMAGE):$(VERSION) --platform=linux/amd64,linux/arm64 ./server/user
	@echo "$(USER_RPC_IMAGE):$(VERSION) built successfully."

push-web:
	@echo "Pushing web image..."
	docker push $(WEB_IMAGE):$(VERSION)
	@echo "$(WEB_IMAGE):$(VERSION) pushed successfully."

push-api:
	@echo "Pushing web image..."
	docker push $(API_IMAGE):$(VERSION)
	@echo "$(API_IMAGE):$(VERSION) pushed successfully."

push-user-rpc:
	@echo "Pushing web image..."
	docker push $(USER_RPC_IMAGE):$(VERSION)
	@echo "$(USER_RPC_IMAGE):$(VERSION) pushed successfully."

# Build all images
build-all: build-web build-api build-user-rpc

# Push all images
push-all: push-web push-api push-user-rpc

build-push-api: build-api push-api
build-push-web: build-web push-web
build-push-user-rpc: build-user-rpc push-user-rpc

# Build and push all images
build-push-all: build-all push-all
	@echo "All Docker images have been built and pushed."