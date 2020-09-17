FROM golang:alpine as build

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /build/
COPY go.mod .
COPY go.sum .
COPY cmd cmd
COPY pkg pkg
COPY main.go .

RUN go mod download
RUN go build -o co-pilot

FROM scratch
COPY --from=build /build/co-pilot /bin/co-pilot
ENTRYPOINT ["/bin/co-pilot"]