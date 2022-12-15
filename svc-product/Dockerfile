FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY main.go .

RUN go build -o http-main .

FROM alpine:3.16 As runner

RUN apk add --no-cache tzdata curl
ENV TZ=Asia/Jakarta

WORKDIR /app

COPY --from=builder /app/http-main .

HEALTHCHECK --start-period=10s --interval=5s --timeout=10s --retries=3 CMD curl -f http://localhost:5001/ping

EXPOSE 5001

CMD ./http-main
