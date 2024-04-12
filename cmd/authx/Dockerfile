FROM builder-cache as builder

LABEL maintainer="kongzz"

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct,https://proxy.golang.com.cn,direct

WORKDIR /server
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o authx-apiserver cmd/authx/main.go

FROM ubuntu:20.04

WORKDIR /server

COPY --from=builder /server/authx-apiserver .

CMD ["/server/authx-apiserver"]
