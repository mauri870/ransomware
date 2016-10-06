#!/bin/bash
set -e

# Build the image
echo "Building the image..."
docker build -t ransomware .

# Compile the binaries
echo "Compiling the binaries..."

if [ $# -eq 0 ]
  then
    echo "Please inform a command to execute inside the container"
    exit
fi

docker run --rm -v "$PWD":/go/src/github.com/mauri870/ransomware ransomware "$@"

# We need change the root permissions of the binaries generated
echo "Fix binaries permissions..."
sudo chown -R $USER:$USER ./bin/
