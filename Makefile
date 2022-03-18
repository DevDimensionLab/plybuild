.DEFAULT_GOAL := all

build:
	go build

docker-build:
	docker build --tag co-pilot:latest .

docker-run:
	docker run co-pilot $(ARGS)

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
	goreleaser --rm-dist

brews-publish:
	goreleaser release --rm-dist --skip-validate -f .goreleaser.brews.yml

upgrade:
	go get github.com/co-pilot-cli/mvn-pom-mutator
	go get -u ./...
	go clean

all: build install


.PHONY: test
