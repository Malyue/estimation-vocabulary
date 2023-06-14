FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

# 复制配置文件
#COPY --from=builder /etc ./etc/

COPY ./etc/config.yaml /etc

EXPOSE 8080

COPY --from=builder /build/main /
ENTRYPOINT ["/main"]