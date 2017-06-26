FROM golang:latest
ENV WORKDIR $GOPATH/src/github.com/mauri870/ransomware
ADD . $WORKDIR
WORKDIR $WORKDIR
RUN go get -u github.com/akavel/rsrc
VOLUME ["$GOPATH/src/github.com/mauri870/ransomware/bin"]