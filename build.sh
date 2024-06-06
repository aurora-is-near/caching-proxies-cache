#!/bin/bash

IMAGE_NAME="caching-proxies-cache"
TAG="latest"

# Enable script to exit on error and print commands and their arguments as they are executed.
set -euxo pipefail

# Check if the first command-line argument is provided and if it's one of the allowed values
if [[ $# -gt 0 ]]; then
    if [[ "$1" == "testnet" || "$1" == "latest" ]]; then
        TAG="$1"
    else
        echo "Error: TAG value must be either 'testnet' or 'latest'."
        exit 1
    fi
fi

echo "Building Docker image '${IMAGE_NAME}:${TAG}'..."
docker build -t "${IMAGE_NAME}:${TAG}" -f Dockerfile .
if [ $? -eq 0 ]; then
    echo "Docker image built successfully."
else
    echo "Docker image build failed."
    exit 1
fi

echo "Process completed successfully."

