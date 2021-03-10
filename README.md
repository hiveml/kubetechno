# `kubetechno`
                                                                             
Kubernetes Tooling for Endpoint Coordination and Host Networking Orchestration

## About 
kubetechno helps workloads run on Kubernetes using host networking.

It does this via the following core goals:

1. Tracking ports on a per-node basis, for the purpose of assigning pods to nodes.
2. Assigning specific ports to pods, in order to prevent port conflicts between pods on the same node.
3. Informing pods of their assigned ports.

Non-core, but still important, goals include:
1. Assisting with health checks.
2. Assisting with service discovery.

See the `docs` directory for more information.

## Is this prod ready?

Not yet, but you can help!
