#!/usr/bin/env bash


docker pull copilotcli/devdimensionlab:latest
docker build --tag co-pilot-test:latest .

docker run -it --rm \
	-v $(pwd)/order:/order \
	co-pilot-test:latest
