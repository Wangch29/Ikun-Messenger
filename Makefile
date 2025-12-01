# define tool installation path
LOCAL_BIN := $(shell pwd)/bin

# export PATH
export PATH := $(LOCAL_BIN):$(PATH)

# define tool versions
PROTOC_GEN_GO_VERSION := latest
PROTOC_GEN_GO_GRPC_VERSION := latest

PROTOC_VERSION := 25.1
PROTOC_ZIP := protoc-$(PROTOC_VERSION)-linux-x86_64.zip
PROTOC_URL := https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)

$(LOCAL_BIN)/protoc:
	@echo "🔍 Protoc not found. Downloading..."
	@mkdir -p $(LOCAL_BIN)
	wget -q -O $(PROTOC_ZIP) $(PROTOC_URL)
	unzip -o $(PROTOC_ZIP) bin/protoc -d .
	unzip -o $(PROTOC_ZIP) "include/*" -d .
	chmod +x $(LOCAL_BIN)/protoc
	rm -f $(PROTOC_ZIP)
	@echo "✅ Protoc installed locally."

# install dependencies to bin/
deps: $(LOCAL_BIN)/protoc
	@mkdir -p $(LOCAL_BIN)
	@echo "Installing dependencies..."
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

all: proto

# Generate Proto
proto: deps
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       api/raftpb/raft.proto

	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       api/kvpb/kv.proto

	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       api/impb/im.proto

clean:
	rm -f api/raftpb/*.pb.go
	rm -f api/kvpb/*.pb.go
	rm -rf bin/

.PHONY: all deps proto clean
