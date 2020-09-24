#!/usr/bin/env bash

set -euo pipefail
shopt -s inherit_errexit

TAG=$(git describe --exact-match --tags)

docker build --tag co-pilot:"$TAG" .

## GITHUB 
#docker tag co-pilot docker.pkg.github.com/co-pilot-cli/co-pilot/co-pilot:"$TAG"
#docker push docker.pkg.github.com/co-pilot-cli/co-pilot/co-pilot:"$TAG"

## DOCKERHUB
docker tag co-pilot:"$TAG" docker.io/copilotcli/co-pilot-cli:"$TAG"
docker push docker.io/copilotcli/co-pilot-cli:"$TAG"
docker tag co-pilot:"$TAG" docker.io/copilotcli/co-pilot-cli:latest
docker push docker.io/copilotcli/co-pilot-cli:latest
