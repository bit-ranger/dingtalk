FROM bitranger/dingtalk:0.1.0

COPY dingtalk /go/src/dingtalk

CMD go run /go/src/dingtalk/app/main.go
