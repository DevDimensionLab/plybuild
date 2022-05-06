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

release-dry:
	goreleaser --rm-dist --skip-publish --skip-validate

brews-publish:
	goreleaser release --rm-dist --skip-validate -f .goreleaser.brews.yml

upgrade:
	go get github.com/devdimensionlab/mvn-pom-mutator
	go get -u ./...
	go clean

all: build install


.PHONY: test
