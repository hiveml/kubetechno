# Flow Overview

## Steps
1. Pod creation request.
2. Mutator service sent review request by API server.
3. Mutator responds to the API server with the patches it determines.
5. The scheduler begins scheduling the pod.
6. The scheduler determines the pod's node through a variety of factors, including the kubetechno port count.
7. The interceptor receives a copy of the bind request for the pod.
8. The interceptor makes changes to the pod, largely port info with the node now known, and then approves the bind.
9. Scheduling finishes up.
