.DEFAULT_GOAL := all

build:
	go build

docker-build:
	docker build --tag ply:latest .

docker-run:
	docker run ply $(ARGS)

docker-publish:
	./docker-publish.sh	

install:
	go install

run:
	go run main.go

test:
	go test -v -cover ./...

lint:
	gofmt -w pkg
	gofmt -w cmd

release:
	goreleaser --clean

release-dry:
	goreleaser --clean --skip-publish --skip-validate

release-brew:
	goreleaser release --clean --skip-validate -f .goreleaser.brews.yml

upgrade:
	go get github.com/devdimensionlab/mvn-pom-mutator
	go get -u ./...
	go clean

all: build install


.PHONY: test
