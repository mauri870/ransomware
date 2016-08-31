ifeq ($(GOOS),"windows")
    BUILD := go build
else
    BUILD := env GOOS=windows go build
endif

default:
	$(BUILD)
	cd server && $(BUILD) && cp ./server.exe ../