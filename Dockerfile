FROM golang:1.8
COPY       micro-auth-worker /bin/micro-auth-worker
ENTRYPOINT ["/bin/micro-user-worker"]
