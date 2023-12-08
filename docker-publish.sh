#!/usr/bin/env bash

set -euo pipefail
shopt -s inherit_errexit

TAG=$(git describe --exact-match --tags)

docker build --tag ply:"$TAG" .

## GITHUB 
#docker tag ply docker.pkg.github.com/ply-cli/ply/ply:"$TAG"
#docker push docker.pkg.github.com/ply-cli/ply/ply:"$TAG"

## DOCKERHUB
docker tag ply:"$TAG" docker.io/copilotcli/ply-cli:"$TAG"
docker tag ply:"$TAG" docker.io/copilotcli/ply-cli:latest
docker push docker.io/copilotcli/ply-cli:"$TAG"
docker push docker.io/copilotcli/ply-cli:latest
