#!/usr/bin/env bash
set -e

image_registry="${1:-localhost:5000}"
image_prefix="${2:-kubetechno/}"
image_tag="${3:-tmp}"

interceptor_img="${image_registry}/${image_prefix}interceptor:${image_tag}"
mutator_img="${image_registry}/${image_prefix}mutator:${image_tag}"
consul_client_img="${image_registry}/${image_prefix}consul-client:${image_tag}"

echo "interceptor image:   ${interceptor_img}"
echo "mutator image:       ${mutator_img}"
echo "consul_client image: ${consul_client_img}"

echo "build interceptor binary/image then push image"
# todo: parallelize these builds while keeping the output clean
cd ../../src/interceptor
./build-binary.sh
./build-push-image.sh "${interceptor_img}"
cd -

echo "build mutator binary/image then push image"
cd ../../src/mutator/server
./build-binary.sh
./build-push-image.sh "${mutator_img}"
cd -

echo "build consulClient binary/image then push image"
cd ../../src/consulClient
./build-binary.sh
./build-push-image.sh "${consul_client_img}"
cd -
