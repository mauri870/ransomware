ifeq ($(GOOS),"windows")
    BUILD := go build
else
    BUILD := env GOOS=windows go build
endif

default:
	mkdir build
	$(BUILD) && mv `ls|grep ransomware` build
	cd server && go build && mv `ls|grep server` ../build