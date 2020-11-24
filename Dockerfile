# Start from golang base image
FROM golang:alpine as builder

# ENV GO111MODULE=on

# Add Maintainer info
LABEL maintainer="Shawn Toubeau <shawntoubeau@gmail.com>"

# Install git
RUN apk update && apk add --no-cache git

# Set wokring directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Downloads all dependencies. Will be cached if go.mod and go.sum don't change
RUN go mod download

# Copy source from current directory to working directory
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start new stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy pre-built binary and env file from previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Run the app!
CMD ["./main"]