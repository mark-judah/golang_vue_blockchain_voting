FROM alpinelinux/golang:latest
USER root
RUN mkdir /app
ADD . /app
WORKDIR /app/raft
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown appuser:appgroup /app
RUN go clean --modcache
RUN go mod download
RUN go build -o raftApp raft
# USER appuser

ENTRYPOINT [ "/app/raft/raftApp" ]
