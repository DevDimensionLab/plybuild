#!/usr/bin/env bash


docker pull plycli/devdimensionlab:latest
docker build --tag ply-test:latest .

docker run -it --rm \
	-v $(pwd)/order:/order \
	ply-test:latest
