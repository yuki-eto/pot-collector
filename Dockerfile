FROM golang:1.10

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/github.com/yuki-eto/pot-collector
ADD . $GOPATH/src/github.com/yuki-eto/pot-collector
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
