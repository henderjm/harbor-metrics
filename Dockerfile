FROM golang:1.12-stretch

COPY main.go /

CMD go run /main.go

