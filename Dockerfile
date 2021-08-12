# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.11
ARG GO_VERSION=1.16.7

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

ENV GOPATH=

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -mod=vendor \
    -installsuffix 'static' \
    -o /app ./cmd/app/main.go

# Final stage: the running container.
FROM scratch AS final

# Import the compiled executable from the first stage.
COPY --from=builder /app /app

# Run the compiled binary.
ENTRYPOINT ["/app"]
