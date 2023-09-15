FROM alpinelinux/golang:latest
USER root
RUN mkdir /app
ADD . /app
WORKDIR /app/vote_backend
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN go clean --modcache
RUN go mod download
RUN go build -o voteApp vote_backend
RUN chown -R appuser:appgroup /app
USER appuser

ENTRYPOINT [ "/app/vote_backend/voteApp" ]
