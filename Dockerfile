FROM golang:1.18 as builder

#
RUN mkdir -p $GOPATH/src/github.com/dilmurodov/online_banking_service 
WORKDIR $GOPATH/src/github.com/dilmurodov/online_banking_service

# Copy the local package files to the container's workspace.
COPY . .

# installing depends and build
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    go mod vendor && \
    make build && \
    mv ./bin/online_banking_service /

FROM alpine
COPY --from=builder online_banking_service .
ENTRYPOINT ["/online_banking_service"]