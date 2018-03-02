FROM golang:1.9
MAINTAINER Mark C Allen <mark@markcallen.com>

RUN wget https://github.com/golang/dep/releases/download/v0.3.1/dep-linux-amd64 -O /usr/local/bin/dep && \
    chmod +x /usr/local/bin/dep

WORKDIR /gopath/src/github.com/markcallen/deezee

ENV GOPATH /gopath
ENV GOBIN /gopath/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

RUN mkdir /deezee

COPY . ./
RUN dep ensure
RUN go build *.go
RUN go install *.go
ENTRYPOINT ["deezee"]
