# Use the official Ubuntu image from the Docker Hub
FROM ubuntu:24.04

# Set the Go version as an environment variable
ENV GO_VERSION=1.20.5

# Install necessary packages
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    git \
    build-essential \
    clang \
    lsb-release && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Download and install Go
RUN wget https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
    rm go$GO_VERSION.linux-amd64.tar.gz

# Set Go paths
ENV PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# Verify Go installation
RUN go version

# Set working directory
WORKDIR /go

# Display Go environment variables
RUN go env

# Default command to run when container starts
# CMD ["bash"]

WORKDIR /app

COPY . .
# Command to run the Go application
CMD ["go", "run", "main.go"]
