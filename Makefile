LOCAL_BIN:=$(CURDIR)/bin
PATH  := $(PATH):$(PWD)/bin

.PHONY: bin-deps
bin-deps:
	$(info Installing binary dependencies...)
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    GOBIN=$(LOCAL_BIN) go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.15

.PHONY: gen-proto-protoc
gen-proto-protoc:
	protoc -I . -I api \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=./pkg --go_opt=paths=source_relative \
    --plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative \
	api/service.proto \
	easyp-demo-service/api/service.proto

.PHONY: generate
generate:
	@$(LOCAL_BIN)/easyp generate

.PHONY: lint
lint:
	@$(LOCAL_BIN)/easyp lint --path api

.PHONY: breaking
breaking:
	@$(LOCAL_BIN)/easyp breaking --against main --path api
