#!/usr/bin/env bash

# This is to start a consul server on one's laptop for local running the churner locally.

bind_ip="${1:-192.168.65.3}" # may depend on os and docker version, this default is a vm address used on a Mac by Docker
docker run -d --net host --name dev-consul consul agent -server -dev -bind="${bind_ip}" -ui -client=0.0.0.0
