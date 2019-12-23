FROM golang:alpine AS builder
WORKDIR /tmp/gopath/src/github.com/wzshiming/console
COPY . .
ENV GOPATH=/tmp/gopath/
ENV GOBIN=$GOPATH/bin/
RUN CGO_ENABLED=0 go install github.com/wzshiming/console/cmd/web_console

FROM alpine
LABEL maintainer="wzshiming@foxmail.com"
COPY --from=builder /tmp/gopath/bin/web_console /usr/local/bin/
VOLUME /var/run/docker.sock
EXPOSE 8888
ENTRYPOINT [ "web_console" ]
