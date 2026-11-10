## Builder image
FROM golang:1.26.10-alpine3.24@sha256:87c65a688975c0b1d4585d204874c45ca78354670e3a5f11368c5a492303d696 AS build

ARG PLY_ACCEPTANCE_RUN_ID=unmanaged
LABEL org.plyproject.acceptance-run=$PLY_ACCEPTANCE_RUN_ID

ENV GO111MODULE=on CGO_ENABLED=0
ENV PLY_GO=/usr/local/go/bin/go

RUN apk --update add bash make git less openssh curl python3 && \
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

COPY --from=build /build/.build/bin/ply /bin/ply
COPY --from=build /build/ktlint /bin/ktlint

ENTRYPOINT ["/bin/ply"]
