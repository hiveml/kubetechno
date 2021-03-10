#!/usr/bin/env bash

# Call as add_extended_resource_to_node.sh <nodeName> <numberOfPorts>. Run kubectl proxy first.
# Each node should have a number of ports assigned equal to the length of the port range.

curl --header "Content-Type: application/json-patch+json" \
     --request PATCH \
     --data "[{\"op\": \"add\", \"path\": \"/status/capacity/kubetechno~1port\", \"value\":\"${2}\"}]" \
     http://localhost:8001/api/v1/nodes/"${1}"/status
