FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY main.go .
RUN go mod init feishu-adapter && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM alpine:3.19
WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 8080

CMD ["./app"]