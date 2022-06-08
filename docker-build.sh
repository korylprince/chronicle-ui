#!/bin/bash

version=$1

tag="korylprince/chronicle-ui:$version"

docker build --no-cache --build-arg "VERSION=$version" --tag "$tag" .

docker push "$tag"

if [ "$2" = "latest" ]; then
    docker tag "$tag:$version" "$tag:latest"
    docker push "$tag:latest"
fi
