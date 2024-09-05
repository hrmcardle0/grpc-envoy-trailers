FROM ubuntu:18.04

# Set working directory
WORKDIR /app

# Install curl with error handling
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy files from the current directory to /app in the container
COPY . .

# Keep the container running indefinitely (for debugging or other purposes)
ENTRYPOINT ["/bin/sh", "-c", "tail -f /dev/null"]
