FROM golang:1.13.1-alpine AS build

WORKDIR /go/src/go.dev.pztrn.name/fastpastebin
COPY . .

WORKDIR /go/src/go.dev.pztrn.name/fastpastebin/cmd/fastpastebin

RUN go build

FROM alpine:3.10
LABEL maintainer "Stanislav N. <pztrn@pztrn.name>"

COPY --from=build /go/src/go.dev.pztrn.name/fastpastebin/cmd/fastpastebin/fastpastebin /app/fastpastebin
COPY docker/fastpastebin.docker.yaml /app/fastpastebin.yaml

EXPOSE 25544
ENTRYPOINT [ "/app/fastpastebin", "-config", "/app/fastpastebin.yaml" ]
