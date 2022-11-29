FROM golang:1.18-alpine AS builder
RUN apk update && apk add --no-cache make git
WORKDIR /go/src/github.com/forbole/bdjuno
COPY . ./
RUN apk add gcc libc-dev
RUN go mod download
RUN make build

FROM alpine:latest as dev
WORKDIR /bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
RUN bdjuno init
COPY ./yml/dev-config.yaml /root/.bdjuno/config.yaml
COPY ./yml/dev-genesis.json ./genesis.json
RUN bdjuno parse genesis-file --genesis-file-path ./genesis.json
CMD ["bdjuno", "start"]

FROM alpine:latest as prod
WORKDIR /bdjuno
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno
RUN bdjuno init
COPY ./yml/prod-config.yaml /root/.bdjuno/config.yaml
COPY ./yml/prod-genesis.json ./genesis.json
RUN bdjuno parse genesis-file --genesis-file-path ./genesis.json
CMD ["bdjuno", "start"]




