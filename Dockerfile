FROM code.pztrn.name/containers/mirror/golang:1.19.0-alpine AS build

WORKDIR /fastpastebin
COPY . .

WORKDIR /fastpastebin/cmd/fastpastebin

RUN CGO_ENABLED=0 go build -tags netgo

FROM code.pztrn.name/containers/mirror/alpine:3.16.1
LABEL maintainer "Stanislav N. <pztrn@pztrn.name>"

COPY --from=build /fastpastebin/cmd/fastpastebin/fastpastebin /app/fastpastebin
COPY docker/fastpastebin.docker.yaml /app/fastpastebin.yaml

EXPOSE 25544
ENTRYPOINT [ "/app/fastpastebin", "-config", "/app/fastpastebin.yaml" ]
