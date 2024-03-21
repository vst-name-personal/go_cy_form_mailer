# Build Stage
FROM golang:1.22 AS build

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/main.go

# Final Stage
FROM alpine:latest

# Meta
LABEL maintainer="***REMOVED***>"
LABEL description="Clean Year Form handler"
LABEL version="0.1"
LABEL release-date="2024-03-20"

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from the build stage
COPY --from=build /app/app .

# Expose ports
EXPOSE 8080

# Run the executable
CMD ./app
