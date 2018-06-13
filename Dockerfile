FROM golang:1.10.3-stretch
RUN go get -u github.com/googlecloudplatform/gcsfuse

FROM debian:stable-slim
COPY --from=0 /go/bin/gcsfuse /usr/local/bin
COPY init.sh /init.sh
