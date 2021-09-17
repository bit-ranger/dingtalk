FROM golang:1.17

COPY dingtalk /go/src/dingtalk
COPY github.com /go/src/github.com
COPY golang.org /go/src/golang.org

ENV GOPATH=/go
ENV GO111MODULE=off

CMD cd /go/src/dingtalk
CMD go build dingtalk/app
