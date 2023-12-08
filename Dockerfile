## Builder image
FROM golang:alpine as build

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk --update add make git less openssh curl && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /build/
COPY . .

RUN go mod download
RUN make test
RUN make build

RUN curl -sSLO https://github.com/pinterest/ktlint/releases/download/0.43.0/ktlint && \
      chmod a+x ktlint

## Shipping image
FROM alpine

RUN apk --update add git less openssh maven graphviz && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

COPY --from=build /build/ply /bin/ply
COPY --from=build /build/ktlint /bin/ktlint

ENTRYPOINT ["/bin/ply"]
