FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git make bash

COPY . /geoipserver

WORKDIR /geoipserver

RUN make server

FROM alpine:latest

WORKDIR /

COPY --from=builder /geoipserver/bin/geoipserver /usr/local/bin/

EXPOSE 8080

VOLUME ["/usr/local/share/GeoIP"]

ENTRYPOINT ["geoipserver"]
