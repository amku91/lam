FROM golang:1.10.3

ADD .  /go/src/github.com/amku91/lam/api
WORKDIR  /go/src/github.com/amku91/lam/api
RUN go get ./...

  # Install api binary globally within container
RUN go install github.com/amku91/lam/api

  # Set binary as entrypoint
ENTRYPOINT /go/bin/api

  # Expose default port (3000)
EXPOSE 8080

