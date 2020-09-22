#!/usr/bin/env bash

set -euo pipefail
shopt -s inherit_errexit

TAG=$(git describe --exact-match --tags)

docker build --tag co-pilot:"$TAG" .
docker tag co-pilot docker.pkg.github.com/co-pilot-cli/co-pilot/co-pilot:"$TAG"
docker push docker.pkg.github.com/co-pilot-cli/co-pilot/co-pilot:"$TAG"
