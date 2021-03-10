# Operating

The project has a [helm](https://helm.sh/) chart as the "standard" way of deploying.

The steps to setting up kubetechno are as follows:

- Set up nodes
  - Run the `set-up-node.sh` script for each node (view script for args).
    Make sure the number of ports added to each node equals the size of the port range. 
  - Using `xargs` and a list of nodes can help automate this.
- Set up namespaces
  - Label all namespaces that will use kubetechno with `kubetechno:user`
- Set up consul, if desired
  - This is not (currently) handled by kubetechno.
  - Make sure appropriate access exists from the host network to consul and vice-versa.
  - Consul has some pretty complex k8s integration options, don't use them for kubetechno.
  - Put at an agent on each node with kubetechno ports.
- Deploy the chart 
  - Create or identify namespace for kubetechno to be deployed to.
  - Customize the helm chart as needed. If there's something you'd like to have the Chart flex on,
    and it's not, feel free to raise an issue.
  - Determine arguments for deploying kubetechno helm chart.
  - Run the install command. 
  - Make sure the pods have started up and are healthy.
- Test the set up
  - If using consul, run the churner example and observe the output of the churner pod.
  - If not using consul try running the manifests under `deploy/examples/simple`.
