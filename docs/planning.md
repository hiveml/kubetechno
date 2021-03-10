# Planning

## Potential Future Features to Add

- Caching port availability *
- A UI
- More/multiple flexible port ranges *
- GC for node locks *
- More documentation, especially for Consul *
- Additional port assignment communication methods
- Versioning
- Explore alternative Consul sync methods

\* High Priority

## Approaches

### Accomplishing Goal 1
Goal 1 can be accomplished by using an 
[extended resource](https://kubernetes.io/docs/tasks/administer-cluster/extended-resource-node/).
A more custom approach would bring little, if any, added benefit.

### Accomplishing Goal 2
One of `kubetechno`'s design guidelines is to limit state and communication to its own components and Kubernetes.
This approach makes state and communication both easier to track and control access to. 
It also keeps kubetechno independent of particular cluster configurations/technology stacks.
A counter-example would be to manipulate containers via their nodes' container run-times.

These state and communication limits are in-part why Goal 2 specifies adding the ports to pods, not containers.
However, the emphasis on pod state also helps link the port assignment to the pod- rather than container-lifecycle.
This lines up better with various components of Kubernetes functionality and keeps the end-points more stable.

Given these motivations, goal 2 is met via an "orchestrator" that communicates with the Kubernetes API.
How and when the orchestrator communicates is limited by goal 3 and is discussed in that goal's section.
However, it's worth noting that the orchestrator is triggered by the API (instead of polling it for example) and 
currently the orchestrator does not keep state between requests as it uses the pods' state from the API.
Given its request based nature and need to communicate to the k8s API, the orchestrator is built into an HTTP server.

In order to secure the state more thoroughly and reduce load on the API, caching the assigned port state is a future 
goal (and has been successfully prototyped).
  
### Accomplishing Goal 3

The design space of how the port assignment is communicated, when in the lifecycle of a Pod the assignment occurs, and 
when that assignment is communicated is complicated not just by the variety of potential ways those 3 specifications can
be meet, but by how those ways relate to and are compatible with one another. 

As such, various possible approaches will be briefly described first, with a conclusion at the end.
While considering the various approaches, it should be noted that they can be altered and combined with each other.

#### Pod to Orchestrator
Perhaps the simplest approach is for a container of a pod to send a request to an orchestrator for ports to request and
then process the response as needed into various configurations for the application components of the pod.

Unfortunately, getting such data into the application is fairly cumbersome. There are several issues, here are some:
- Either a pod needs its app containers to handle kubetechno integration before starting or use an init container.
- Both approaches would require extra/custom efforts to be made to get configuration into the apps.
- Multi-container pods may need to share the configuration.
- A need for custom configuration makes 3rd party images especially hard to use.
- The port assignments are difficult to get as env variables, especially for use in container command args.

#### Assigned Node Scheduling Approach
A mutator assigns a node to a pod using a node selector and then assigns available ports on the now fixed node. 
However, avoiding the scheduler is extremely undesirable, so this approach is as well. 

#### Kubelet Device Plugin Approach
By treating host ports as devices, a device plugin would allow env vars and volumes to be mounted to containers. 
However, little info is passed to the plugin making it hard to track port assignments.
Thus, it would be difficult to determine when ports should be freed and made assignable again. 

Also, a device plugin would not share the port information with all of a pod's containers (without volume mounts).
This would make multi-container pods more cumbersome to properly configure. 

Furthermore, an extra step would be required to save state to k8s in order to make the info available via the k8s API.

#### Request Resource Alternation Approach
An init container could request ports and then the orchestrator communicates the assignment info by changing resources, 
such as pod attributes, that can be accessed, such as via field references, by the pod to serve as env variables. 
Once that info is updated on the pod's node's Kubelet's cache, the info can be used by the pod's non-init containers. 

There are two fairly large problems with this.
First, it relies on the kubelet cache updating which would require periodic polling of Kubelet or a long wait period. 
Second, it is not a particularly reliable mechanism to base such a system on as Kublet could be redesigned so that the 
info that is determined at pod startup is locked in and not refreshed as before the non-init containers are spun up.

#### Missing ConfigMap Approach

A pod can be scheduled to a node while referencing a config map that does not exist.
The pod will then enter an error state, but on a start-up retry if the configmap has been added the pod will start.
It should be noted that the config map has to be missing completely, not just missing a key being referenced.
Pods in the error state can be found using an API watch or some other mechanism, such as polling.
An important implementation component is that the pod's would need to reference a unique config-map.
To aid in this a mutation admission web-hook can be employed to alter the pods' config map references on pod creation.

One benefit to this approach is that config-maps are quite flexible in how their information is passed into pods.

There are some substantial downsides with this approach though:
- Many config-maps need to be generated.
- The pods have to enter an error state.
- The pods in the error state have to be queried for, likely through an API stream/watch or periodic polling.
- This exploit may not work with future releases. Kubernetes may change and not schedule pods with missing config-maps.

### Altered ConfigMaps

This is somewhat of a combination of the two previous approaches.
Instead of requesting annotation changes like the "Request Resource Alternation Approach", a config-map is altered.
This approach's pros and cons are a mixture of the previous two approaches.

#### Scheduling Plugin or Extender

The scheduler allows for "interceptions" to occur at particular points in the scheduling timeline.
Of particular interest is after the pod's node has been decided upon, but before the pod is bound to the node.
This step is "pre-bind" for a scheduler plugin and can be approximated via an extender via the "bind" step by simply 
binding the non-kubetechno pods and running the pre-bind logic, before subsequently binding any kubetechno pods.

Before binding, the logic for accomplishing goal 2 can be invoked and ports to assign determined.
Then, looking back to the other approaches, this assigment can be communicated, likely in one of the following ways:

- Via configmaps set up similarly to the "Altered ConfigMaps" approach.
- Via pod changes that can be used in field references.
- Via an endpoint that can be called by the containers. 

Potentially multiple communication methods could be used simultaneously or as selected by a cluster's administrators. 

This approach's downsides are typical of custom scheduling sub-systems and the downsides of the communication method(s).

A plugin seems to be more preferable to an extender with no network communication or pod binding requirements.
However, the plugin's seems more cumbersome to install and manage.
The plugin also has to consider pods that are not bound but have passed through the plugin goal 2.
The extender can hold a node lock until a pod is processed and then bound (admittedly this might be slower).
A compromise exists in a plugin that behaves like the extender but operates at the bind phase of the scheduler.

#### Bind Interceptor 

The "Scheduling Plugin or Extender" approach picks likely the most optimal timing of "when" goal 2 should be met.
However, it may not be ideal "how" as it is a change to the k8s scheduler (or requires an additional scheduler).
A Validating Admission Controller for bind create requests to the API server intercepts at the same point.
Goal 2 logic could be done after when the admission controller is called, and before it approves the bind creation.

Webhook Admission Controllers have a very clear configuration schema, clean operations, and secure communications.
These traits plus being more independent of the scheduler seem to make this a more preferable option.
The approach also leaves open the same communication approaches as a scheduling plugin or extender.

Notes: 
- A validating admission controller is used instead of a mutating one as the bind itself doesn't need to be changed.
- Something worth exploring in the future for the server to bind the pod and reject the bind request. This would remove
the 'after hitting the interceptor but before bind' stage, make tracking the state of pods simpler. 

#### Approach Selection

After exploring the various options the Bind Admission controller seems best to trigger the Goal 2 logic.

That said some options stand out as being potentially worth revisiting in the future if the requirements of kubetechno 
or the mechanisms of Kubernetes change in ways favorable to their usage, or if the analysis needs to be revisited. 
The selection of the approach to Goal 3 is just a means to an end, so please raise an issue if there is an alternative 
approach that may be better suited (of your own devising or one listed above).

The use of a bind admission controller doesn't limit the communication method to the pod very much.
For now, changing the pods' annotations hits a sweet sport of lower operational overhead with ease of use/power in 
manifests. In the future additional communication methods may be added. 

### Non-Core Goals Approaches

The non-core goals are accomplished via integrating the pod's readiness probes with [Consul](https://www.consul.io/). 
There several ways to do so and a more in-depth discussion is forthcoming.
The mutator is designed such that additional integrations could be merged in fairly easily.
