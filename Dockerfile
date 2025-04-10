# Use Go 1.23 bookworm as base image
FROM golang:1.23-bookworm AS base

# Set working directory inside the container
WORKDIR /golangProject

# Copy go.mod and go.sum to the /app directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . /golangProject/

# List files in /app for debugging
RUN ls -al /golangProject

# Change to the directory where the main.go file is located
WORKDIR /golangProject/cmd/app

# List files in /app/cmd/app for debugging
RUN ls -al /golangProject/cmd/app

# Build the application
RUN go build -o /main

# Expose the port
EXPOSE 8080

# Command to run the application
CMD ["/main"]
