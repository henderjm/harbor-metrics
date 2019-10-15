FROM golang:1.12-stretch

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o harbor-metrics

CMD ["./harbor-metrics"]

