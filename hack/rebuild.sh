#!/usr/bin/env bash

set -e

GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o test/devpod-cli-linux-amd64
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o test/devpod-cli-linux-arm64
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o test/devpod-cli-darwin-arm64
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o test/devpod-cli-darwin-amd64
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o test/devpod-cli-windows-amd64
sudo go build -o /usr/local/bin/devpod
cp test/devpod-cli-linux-amd64 test/devpod-linux-amd64
cp test/devpod-cli-linux-arm64 test/devpod-linux-arm64
cp test/devpod-cli-linux-amd64 desktop/src-tauri/bin/devpod-cli-x86_64-unknown-linux-gnu
cp test/devpod-cli-linux-arm64 desktop/src-tauri/bin/devpod-cli-aarch64-unknown-linux-gnu
cp test/devpod-cli-darwin-amd64 desktop/src-tauri/bin/devpod-cli-x86_64-apple-darwin
cp test/devpod-cli-darwin-arm64 desktop/src-tauri/bin/devpod-cli-aarch64-apple-darwin


#upx "test/devpod-linux-amd64"
#upx "test/devpod-linux-arm64"
#upx "test/devpod-darwin-arm64"
rm -R $TMPDIR/devpod-cache 2>/dev/null || true
