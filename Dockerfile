# Use the official Go image as a base
FROM --platform=linux/amd64 golang:1.22-alpine

# Install necessary tools
RUN apk update && apk add --no-cache mysql-client

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Set environment variables for cross-compilation
ENV GOOS=linux GOARCH=amd64

# Build the Go application
RUN go build -o db-recordum main.go

# Command to run the application
CMD ["/app/db-recordum"]
