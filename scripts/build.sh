#!/bin/sh

# simple build script, does the job

echo "installing go packages..."

go get

echo "building binary..."

go build -o server .
chmod 755 server

echo "running server..."
./server