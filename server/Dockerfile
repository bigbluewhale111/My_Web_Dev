# Use the official Golang image as a base
FROM golang:1.22.4

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application to the container
COPY . .

# Build the Go application
RUN go build main.go

# Expose port 3000 (or whatever port your Go server uses)
EXPOSE 3000

# Command to run the executable
CMD ["./main"]