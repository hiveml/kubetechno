#!/usr/bin/env bash
set -e

img="${1:-locahost:5000/kubetechno/consul-client:tmp}"

# build
docker build -t "${img}" .
echo ""
echo "docker img built: ${img}"
echo ""

# push
docker push "${img}"
echo ""
echo "docker img pushed: ${img}"
echo ""
