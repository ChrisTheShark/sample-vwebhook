FROM golang:1.14-alpine3.12 AS build

RUN mkdir -p /go/src/github.com/ChrisTheShark/simple-vwebhook
WORKDIR /go/src/github.com/ChrisTheShark/simple-vwebhook

COPY . .
RUN apk --no-cache add build-base gcc && \
    adduser -S 10001 golang && \
    go test -mod=vendor ./... && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -mod vendor -o main

FROM alpine:3.12

COPY --from=build /go/src/github.com/ChrisTheShark/simple-vwebhook/main .
COPY --from=build /etc/passwd /etc/passwd

USER 10001

ENTRYPOINT [ "/main" ]