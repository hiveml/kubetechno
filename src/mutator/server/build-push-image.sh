#!/usr/bin/env bash

set -e

img="${1:-locahost:5000/kubetechno/mutator:tmp}"

# build
docker build -t "${img}" .
echo ""
echo "docker image built: ${img}"
echo ""

# push
docker push "${img}"
echo ""
echo "docker image pushed: ${img}"
echo ""
