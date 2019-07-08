FROM golang:1.11-alpine

WORKDIR /app

RUN apk add git && \
    go get -u -v gitlab.com/pztrn/fastpastebin/cmd/fastpastebin

FROM alpine:3.10
LABEL maintainer "Stanislav N. <pztrn@pztrn.name>"

COPY --from=0 /go/bin/fastpastebin /app/fastpastebin
COPY examples/fastpastebin.yaml.docker /app/fastpastebin.yaml

EXPOSE 25544
ENTRYPOINT [ "/app/fastpastebin", "-config", "/app/fastpastebin.yaml" ]
