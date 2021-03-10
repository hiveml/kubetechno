# Churner

This directory holds the code and helm chart for running an example demonstrating kubetechno's Consul integration.
On a loop:
- The app pods, service instances, and data returned from the app points is gathered.
- The results of each data source are then compared to ensure everything is aligned.
- Each pod then has a chance of being deleted from Kubernetes (and then replaced), simulating normal pod "churn".

## Operation

### Set Up
- Pick a set of nodes.
- Set up consul on those nodes.
- Add custom port resources to the nodes (the number of ports should equal the port range set for kubetechno).
- Tag at least one node with `kubetechnoTest:user`. 
  That's what where the churner pod will run.
- Run `deploy.sh`, considering the image options. 
  See `deploy.sh` for details.
  
### Observation
Watch the log of the churner pod. 

### Shutdown
To shutdown. just remove the deployed chart.

Optionally:
- Save the logs.
- Remove the consul instances.
- Remove the `kubetechnoTest:user` from nodes.
- Remove the kubetechno port resources from nodes.
