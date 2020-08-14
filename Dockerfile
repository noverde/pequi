# Use the official Golang image to create a build artifact.
FROM golang:1.14 as builder

# Create and change to the app directory.
WORKDIR $GOPATH/src/app

# Copy the code from the host and compile it.
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go mod vendor && \
    go build -a -v -o pequi .

# Use the official Alpine image for a lean production container.
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/app/pequi /app/pequi

# Run the web service on container startup.
CMD ["/app/pequi"]
