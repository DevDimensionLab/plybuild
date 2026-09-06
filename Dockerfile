## Builder image
FROM golang:1.26.7-alpine3.24@sha256:28d89ee9cc0ff9fec75c82ca201e6bf7fdf9a679d4b7b24dfa04f2bb766bb468 AS build

ARG PLY_ACCEPTANCE_RUN_ID=unmanaged
LABEL org.plyproject.acceptance-run=$PLY_ACCEPTANCE_RUN_ID

ENV GO111MODULE=on CGO_ENABLED=0

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

COPY --from=build /build/ply /bin/ply
COPY --from=build /build/ktlint /bin/ktlint

ENTRYPOINT ["/bin/ply"]
