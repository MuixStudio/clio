.PHONY: proto

PROTO_FILES := $(wildcard proto/*.proto)
OUT_DIR := ../common

proto:
	protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_FILES)
