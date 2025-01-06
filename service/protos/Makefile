PROTO_DIR := proto
PROTO_FILE := proto/*/*.proto

GEN_GO_DIR := gen/go
GEN_PYTHON_DIR := gen/python

install:
	@echo "Installing grpc proto generators..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@sudo apt install python3-grpc-tools python3-grpcio -y

go-generate:
	@echo "Generating Go code from $(PROTO_FILE)..."
	@protoc -I $(PROTO_DIR) $(PROTO_FILE) --go_out=$(GEN_GO_DIR) --go_opt=paths=source_relative --go-grpc_out=$(GEN_GO_DIR) --go-grpc_opt=paths=source_relative

python-generate:
	@echo "Generating Python code from $(PROTO_FILE)..."
	@python3 -m grpc_tools.protoc -I=$(PROTO_DIR) --python_out=$(GEN_PYTHON_DIR) --grpc_python_out=$(GEN_PYTHON_DIR) --pyi_out=$(GEN_PYTHON_DIR) $(PROTO_FILE)

all: install go-generate python-generate