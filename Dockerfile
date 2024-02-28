FROM golang:latest

LABEL version="1.0"

RUN mkdir /go/src/caching-proxies-cache
COPY . /go/src/caching-proxies-cache
WORKDIR /go/src/caching-proxies-cache

RUN go mod tidy

ENTRYPOINT ["./entrypoint.sh"]
