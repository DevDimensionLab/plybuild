## Builder image
FROM golang:alpine as build

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk --update add make git less openssh maven && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /build/
COPY . .

RUN go mod download
RUN make test
RUN make build


## Shipping image
FROM alpine

RUN apk --update add git less openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

COPY --from=build /build/co-pilot /bin/co-pilot

CMD ["/bin/co-pilot"]
