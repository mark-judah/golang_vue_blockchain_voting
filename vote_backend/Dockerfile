# FROM alpinelinux/golang:latest
FROM golang:1.21-rc-alpine3.17
ENV CGO_ENABLED=1
# USER root
RUN apk add --update npm
RUN apk add --update build-base
RUN apk add --update sqlite
RUN mkdir /app
WORKDIR /app
ADD . /app
# RUN addgroup -S appgroup && adduser -S appuser -G appgroup
# RUN go clean --modcache
# RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

# RUN chown -R appuser:appgroup /app
# USER appuser

ENTRYPOINT CompileDaemon -build="go build -o voteApp" -command="./voteApp"
