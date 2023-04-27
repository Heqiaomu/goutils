#!/usr/bin/env bash

echo "====> run log check linter"
if ! ./script/linter-logcheck.sh "$(pwd)"; then
  exit 1
fi

echo "====> add git hooks config file"
./script/changelog.sh

echo "====> add private go modules..."
#go get -insecure github.com/Heqiaomu/glog@v1.0.1
#go get -insecure github.com/Heqiaomu/protocol@v1.0.0

echo "====> go mod download..."
go mod tidy
