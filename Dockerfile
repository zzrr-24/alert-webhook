FROM golang:1.25-alpine AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o alert-webhook ./cmd/webhook

FROM alpine:3.19
WORKDIR /root/
COPY --from=builder /app/alert-webhook .
COPY --from=builder /app/configs/config.yaml ./configs/config.yaml

EXPOSE 8080

CMD ["./alert-webhook"]
