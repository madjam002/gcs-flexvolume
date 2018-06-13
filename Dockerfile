FROM golang:1.10.3-stretch
WORKDIR /go/src/github.com/awprice/gcs-flexvolume
ADD Makefile .
ADD main.go .
RUN go get -u github.com/googlecloudplatform/gcsfuse
RUN make build

FROM debian:stable-slim
COPY --from=0 /go/bin/gcsfuse /usr/local/bin
COPY --from=0 /go/src/github.com/awprice/gcs-flexvolume/gcsfuse-driver /usr/local/bin
COPY init.sh /init.sh
