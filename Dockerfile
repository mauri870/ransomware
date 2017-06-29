FROM golang:latest
ENV WORKDIR $GOPATH/src/github.com/mauri870/ransomware
ENV GOROOT /usr/local/go
ADD . $WORKDIR
WORKDIR $WORKDIR
RUN make deps
VOLUME ["$GOPATH/src/github.com/mauri870/ransomware/bin"]