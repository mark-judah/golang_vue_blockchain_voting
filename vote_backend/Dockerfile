FROM alpinelinux/golang:latest
USER root
RUN apk add --update npm
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN go clean --modcache
RUN go mod download
RUN go build -o voteApp 
RUN chown -R appuser:appgroup /app
USER appuser

ENTRYPOINT [ "/app/voteApp" ]
