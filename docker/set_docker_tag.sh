#!/usr/bin/env bash

if [[ $CI_BUILD_REF_NAME == "master" ]]; then
    export DOCKER_TAG=latest;
else
    export DOCKER_TAG="${CI_BUILD_REF_NAME}";
fi