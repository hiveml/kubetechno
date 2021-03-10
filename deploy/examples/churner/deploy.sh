#!/usr/bin/env bash
set -e

image_registry="${1:-localhost:5000}"
image_prefix="${2:-kubetechno/examples/}"
image_tag="${3:-tmp}"
namespace="${4:-default}"

app_img="${image_registry}/${image_prefix}app:${image_tag}"
churner_img="${image_registry}/${image_prefix}churner:${image_tag}"

echo "now installing kubetechno example with the following values:"
echo "  namespace: ${namespace}"
echo "  app_img:   ${app_img}"
echo "  churner_img: ${churner_img}"

helm uninstall --namespace "${namespace}" kubetechno-churner kubetechno-churner || echo "nothing to uninstall"
# the helm install doesn't wait for the pod to be deleted
kubectl delete po -l app=kubetechno-example -n "${namespace}" || echo "no app pods to delete"
kubectl delete po kubetechno-churner || echo "kubetechno churner already deleted"

# to do: parallelize these builds while keeping the output clean
cd src/app
./build-binary.sh
./build-push-image.sh "${app_img}"
cd -

cd src/churner
./build-binary.sh
./build-push-image.sh "${churner_img}"
cd -

helm install --namespace "${namespace}" --set img.app="${app_img}" --set img.churner="${churner_img}" kubetechno-churner kubetechno-churner
