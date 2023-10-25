# Start from Golang 1.19 base image
FROM golang:1.19-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

### Start a new stage from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
