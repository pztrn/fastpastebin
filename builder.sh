#!/bin/bash

VERSION=$1
if [ "${VERSION}" == "" ]; then
    echo "Specify version as first parameter!"
    exit 1
fi

OS_LIST=("darwin/amd64" "dragonfly/amd64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/ppc64" "linux/ppc64le" "linux/mips" "linux/mipsle" "linux/mips64" "linux/mips64le" "linux/s390x" "netbsd/386" "netbsd/amd64" "netbsd/arm" "openbsd/386" "openbsd/amd64" "openbsd/arm" "solaris/amd64" "windows/386" "windows/amd64")

if [ ! -d ./dist ]; then
    mkdir -p ./dist
fi

for os in ${OS_LIST[@]}; do
    mkdir -p ./dist/${os}
    goos=$(echo ${os} | awk -F"/" '{ print $1 }')
    goarch=$(echo ${os} | awk -F"/" '{ print $2 }')
    echo "Building for ${goos} ${goarch}..."
    GOOS=${goos} GOARCH=${goarch} GOFLAGS="-mod=vendor" go build -o ./dist/${os}/fastpastebin ./cmd/fastpastebin/
    cp ./examples/fastpastebin.yaml.dist ./dist/${os}/fastpastebin.yaml
    cd ./dist/${os}/
    tar -czf fastpastebin-${VERSION}-${goos}-${goarch}.tar.gz fastpastebin fastpastebin.yaml
    mv fastpastebin-${VERSION}-${goos}-${goarch}.tar.gz ../../
    cd - &>/dev/null
done