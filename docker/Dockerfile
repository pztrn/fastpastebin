FROM golang:1.11-alpine AS build

WORKDIR /app

RUN apk add git && \
    go get -u -v gitlab.com/pztrn/fastpastebin/cmd/fastpastebin

FROM alpine:3.10
LABEL maintainer "Stanislav N. <pztrn@pztrn.name>"

COPY --from=build /go/bin/fastpastebin /app/fastpastebin
COPY fastpastebin.docker.yaml /app/fastpastebin.yaml

EXPOSE 25544
ENTRYPOINT [ "/app/fastpastebin", "-config", "/app/fastpastebin.yaml" ]
