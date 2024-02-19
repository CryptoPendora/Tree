#Dockerfile

# Use Ubuntu 22.04 as the base image
FROM ubuntu:22.04 AS builder

# Install necessary packages for building your application
RUN apt-get update && apt-get install -y \
    gcc \
    libsqlite3-dev \
    git \
    wget \
    tar \
    ca-certificates \
    nano \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get upgrade    

# Manually install Go 1.22.0 with a retry mechanism
WORKDIR /tmp
RUN wget -c https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Set environment variables for Go
ENV PATH=$PATH:/usr/local/go/bin \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Copy the entire directories for Server and Reader
COPY ./ServerHeartbeatWarriors /app/ServerHeartbeatWarriors
COPY ./ReaderHeartbeatWarriors /app/ReaderHeartbeatWarriors

# Check Go version
RUN go version    

# Download dependencies and build the ServerHeartbeatWarriors application
WORKDIR /app/ServerHeartbeatWarriors
RUN go mod download
RUN go build -tags 'libsqlite3' -o ServerHeartbeatWarriors .

# Download dependencies and build the ReaderHeartbeatWarriors application
WORKDIR /app/ReaderHeartbeatWarriors
RUN go mod download
RUN go build -tags 'libsqlite3' -o ReaderHeartbeatWarriors .

# Use Ubuntu 22.04 as the runtime base image for the final stage
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libsqlite3-0 \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# # Copy the compiled applications from the builder stage
# COPY --from=builder /app/ServerHeartbeatWarriors/ServerHeartbeatWarriors /app/ServerHeartbeatWarriors
# COPY --from=builder /app/ReaderHeartbeatWarriors/ReaderHeartbeatWarriors /app/ReaderHeartbeatWarriors

# สำหรับ builder stage
COPY --from=builder /app/ServerHeartbeatWarriors/ServerHeartbeatWarriors /app/
COPY --from=builder /app/ReaderHeartbeatWarriors/ReaderHeartbeatWarriors /app/

# Set the working directory
WORKDIR /app

# Commands to run the applications (This will be overridden by docker-compose)
CMD ["bash"]

