FROM golang:alpine

############################
# STEP 1 ca-certificates: https://github.com/google/go-github/issues/1049
############################
RUN apk update \
        && apk upgrade \
        && apk add --no-cache ca-certificates \
        && update-ca-certificates 2>/dev/null || true

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Run test
RUN go test ./...

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

COPY ./config.json /dist/config.json
COPY ./index.html /dist/index.html

ENV GITHUB_AUTH_TOKEN=a92bfa63768085cdc3c689f671794fbd9c82e933

# Command to run the executable
ENTRYPOINT ["/dist/main"]