#!/usr/bin/env bash

docker build \
  -f docker/Dockerfile \
  --no-cache \
  --build-arg gitCommitHash=$(git rev-parse HEAD) \
  -t cremindes/whalelint:latest \
  .

