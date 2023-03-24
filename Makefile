BINARY_SERVER=srv
BINARY_CLIENT=cli

SERVER_WINDOWS=$(BINARY_SERVER)_windows_amd64.exe
CLIENT_WINDOWS=$(BINARY_CLIENT)_windows_amd64.exe
SERVER_LINUX=$(BINARY_SERVER)_linux_amd64
CLIENT_LINUX=$(BINARY_CLIENT)_linux_amd64
SERVER_DARWIN=$(BINARY_SERVER)_dawin_amd64
CLIENT_DARWIN=$(BINARY_CLIENT)_dawin_amd64
VERSION=$(shell git describe --tags --always --long --dirty)


server_windows: $(SERVER_WINDOWS)
client_windows: $(CLIENT_WINDOWS)
server_linux: $(SERVER_LINUX)
client_linux: $(CLIENT_LINUX)
server_darwin: $(SERVER_DARWIN)
client_darwin: $(CLIENT_DARWIN)


$(SERVER_WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -v -o $(SERVER_WINDOWS) -ldflags="-s -w -X main.buildVersion=$(VERSION)" ./server/cmd/main.go
$(CLIENT_WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -v -o $(CLIENT_WINDOWS) -ldflags="-s -w -X main.buildVersion=$(VERSION)" ./client/cmd/main.go
$(SERVER_LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(SERVER_LINUX) -ldflags="-s -w -X main.buildVersion=$(VERSION)" ./server/cmd/main.go
$(CLIENT_LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(CLIENT_LINUX) -ldflags="-s -w -X main.buildVersion=$(VERSION)" ./client/cmd/main.go
$(SERVER_DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -v -o $(SERVER_DARWIN) -ldflags="-s -w -X main.buildVersion=$(VERSION)" ./server/cmd/main.go
$(CLIENT_DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -v -o $(CLIENT_DARWIN) -ldflags="-s -w -X main.buildVersion=$(VERSION)" ./client/cmd/main.go



build-all: server_linux client_linux server_windows client_windows server_darwin client_darwin
	@echo version: $(VERSION)
build-linux: server_linux client_linux
	@echo version: $(VERSION)
build-windows: server_windows client_windows
	@echo version: $(VERSION)
build-darwin: server_darwin client_darwin
	@echo version: $(VERSION)



run_server: build-linux
	./$(SERVER_LINUX)
run_client: build-linux
	./$(CLIENT_LINUX)

test:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out

clean:
	go clean
	rm $(SERVER_LINUX)
	rm $(CLIENT_LINUX)
clean-all:
	go clean
	rm $(SERVER_LINUX)
	rm $(CLIENT_LINUX)
	rm $(SERVER_WINDOWS)
	rm $(CLIENT_WINDOWS)
	rm $(SERVER_DARWIN)
	rm $(CLIENT_DARWIN)

deps:
	go mod tidy