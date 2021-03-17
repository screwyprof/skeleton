# Start from golang base image
FROM golang:1.17 as builder

# Add Maintainer info
LABEL maintainer="Maksim Shcherbo <max@happygopher.nl>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git make

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
#COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
#RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN BINARY=main make build-docker

# Start a new stage from scratch
FROM alpine:3.14
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /app/bin/main .

# Expose port to the outside world
EXPOSE 8000

#Command to run the executable
CMD [ "./main" ]