FROM golang:latest
RUN mkdir -p $GOPATH/src/github.com/mauri870/ransomware
ADD . $GOPATH/src/github.com/mauri870/ransomware
WORKDIR $GOPATH/src/github.com/mauri870/ransomware
RUN go get -u github.com/akavel/rsrc
VOLUME ["$GOPATH/src/github.com/mauri870/ransomware"]