FROM golang:1.13.1-alpine AS build

WORKDIR /fastpastebin
COPY . .

WORKDIR /fastpastebin/cmd/fastpastebin

RUN GOFLAGS="-mod=vendor" go build

FROM alpine:3.10
LABEL maintainer "Stanislav N. <pztrn@pztrn.name>"

COPY --from=build /fastpastebin/cmd/fastpastebin/fastpastebin /app/fastpastebin
COPY docker/fastpastebin.docker.yaml /app/fastpastebin.yaml

EXPOSE 25544
ENTRYPOINT [ "/app/fastpastebin", "-config", "/app/fastpastebin.yaml" ]
