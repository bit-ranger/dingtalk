FROM bitranger/dingtalk:latest

COPY dingtalk /go/src/dingtalk

CMD go run /go/src/dingtalk/app/main.go
