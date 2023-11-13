# FROM alpinelinux/golang:latest
FROM golang:1.21-rc-alpine3.17
ENV CGO_ENABLED=1
RUN apk add --update npm
RUN apk add --update build-base
RUN apk add --update sqlite
RUN mkdir /app
WORKDIR /app
ADD . /app
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -build="go build -o voteApp" -command="./voteApp"
