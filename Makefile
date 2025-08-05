PROTO_DIR := proto
PROTO_SRC := $(shell find $(PROTO_DIR) -name "*.proto")
GO_OUT := .

.PHONY: generate-proto
generate-proto:
	@for file in $(PROTO_FILES); do \
		echo "Generating $$file..."; \
		protoc \
			--proto_path=$(PROTO_DIR) \
			--go_out=$(GO_OUT) \
			--go-grpc_out=$(GO_OUT) \
			$$file; \
	done