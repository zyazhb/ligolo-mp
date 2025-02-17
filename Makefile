GO_VER = 1.23.5
BLOAT_FILES = AUTHORS CONTRIBUTORS PATENTS VERSION favicon.ico robots.txt SECURITY.md CONTRIBUTING.md LICENSE README.md ./doc ./test ./api ./misc
GARBLE_VER = 1.23.5

.PHONY: all
all: assets binaries

.PHONY: assets
assets: go agent

.PHONY: binaries
binaries: server client

.PHONY: go
go:
	# Build go
	cd artifacts && curl -L --output go$(GO_VER).linux-amd64.tar.gz https://dl.google.com/go/go$(GO_VER).linux-amd64.tar.gz
	cd artifacts && tar xvf go$(GO_VER).linux-amd64.tar.gz
	cd artifacts/go && rm -rf $(BLOAT_FILES)
	rm -f artifacts/go/pkg/tool/linux_amd64/doc
	rm -f artifacts/go/pkg/tool/linux_amd64/tour
	rm -f artifacts/go/pkg/tool/linux_amd64/test2json
	# Build garble
	cd artifacts/go/bin && curl -L --output garble https://github.com/ttpreport/garble/releases/download/v$(GARBLE_VER)/garble_linux_amd64 && chmod +x garble
	# Bundle
	cd artifacts && zip -r go.zip ./go
	# Clean up
	cd artifacts && rm -rf go go$(GO_VER).linux-amd64.tar.gz

.PHONY: agent
agent:
	cd artifacts/agent && go mod tidy
	cd artifacts/agent && go mod vendor
	cd artifacts/agent && zip -r ../agent.zip .

.PHONY: server
server:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -trimpath -o ligolo-server ./cmd/server/

.PHONY: client
client:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -trimpath -o ligolo-client_linux ./cmd/client/
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -trimpath -o ligolo-client_windows ./cmd/client/

.PHONY: protobuf
protobuf:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protobuf/ligolo.proto

.PHONY: clean
clean:
	rm -rf artifacts/agent.zip artifacts/go.zip

.PHONY: install
install:
	./install.sh
