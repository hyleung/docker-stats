FROM golang:1.7.4-alpine

ARG GOOS
ARG GOARCH

ENV GOOS $GOOS
ENV GOARCH $GOARCH

WORKDIR /go/src/github.com/hyleung/docker-stats
COPY . ./

RUN go get