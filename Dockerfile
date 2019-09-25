FROM golang:1.13-alpine AS build

WORKDIR /go/src/gitlab.com/pztrn/fastpastebin
COPY . .

WORKDIR /go/src/gitlab.com/pztrn/fastpastebin/cmd/fastpastebin

RUN go build

FROM alpine:3.10
LABEL maintainer "Stanislav N. <pztrn@pztrn.name>"

COPY --from=build /go/src/gitlab.com/pztrn/fastpastebin/cmd/fastpastebin/fastpastebin /app/fastpastebin
COPY docker/fastpastebin.docker.yaml /app/fastpastebin.yaml

EXPOSE 25544
ENTRYPOINT [ "/app/fastpastebin", "-config", "/app/fastpastebin.yaml" ]
