# Start from golang base image
FROM golang:1.18-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git make

# add user
RUN adduser -D gouser && chown -R gouser /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN BINARY=main make build-docker

# Start a new stage from scratch
FROM scratch

# Add Maintainer info
LABEL maintainer="Maksim Shcherbo <max@happygopher.nl>"

WORKDIR /root/

# copy ca certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy users from builder
COPY --from=builder /etc/passwd /etc/passwd

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /app/bin/main .

# set user as as the owner of the app
USER gouser

# Expose port to the outside world
EXPOSE 8000

#Command to run the executable
CMD [ "./main" ]