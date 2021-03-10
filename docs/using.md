# Using

Adding the annotation `kubetechno_port_count: "<a postive int>"` to a pod or pod template signals to the kubetechno machinery 
to make appropriate changes to the pod(s) (it does not make changes to pod templates). 

To illustrate the changes that are made, what follows is an example deployment manifest and subsequently simplified pod 
information for a pod created with that template using kubetechno and kubetechno-Consul integration. 
The information is simplified by the info not relevant to kubetechno being removed.

### Example Deployment Yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubetechno-churner-app
  namespace: default
spec:
  selector:
    matchLabels:
      app: kubetechno-churner-app
  replicas: 2
  template:
    metadata:
      annotations:
        kubetechno_port_count: "{{ .Values.portCount }}"
        kubetechno_consul_service: kubetechno-churner-app
        kubetechno_consul_service_pull_policy: Always
        kubetechno_consul_check_path: "/"
        kubetechno_consul_initial_delay_seconds: "3"
        kubetechno_consul_period_seconds: "30"
      labels:
        app: kubetechno-churner-app
        kubetechno: user
    spec:
      containers:
        - name: server
          command:
            - ./server 
          args:
            - ${PORT0} 
          image: localhost:5000/kubetechno/example-server
        - name: sidecar
          image: localhost:5000/kubetechno/example-sidecar
```

### Simplified Pod Yaml
```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    PORT0: "9000"
    PORT1: "9001"
    kubetechno: active
    kubetechno-node: docker-desktop
    kubetechno_consul_check_path: /
    kubetechno_consul_initial_delay_seconds: "3"
    kubetechno_consul_period_seconds: "30"
    kubetechno_consul_service: kubetechno-churner-app
    kubetechno_consul_service_pull_policy: Always
    kubetechno_port_count: "2"
  labels:
    app: kubetechno-churner-app
    kubetechno: user
    kubetechno-node: docker-desktop 
  name: kubetechno-churner-app-5fd5988c4d-tppwx
  namespace: default
spec:
   initContainers:
    - env:
      - name: PORT0
        valueFrom:
          fieldRef:
            apiVersion: v1
            fieldPath: metadata.annotations['PORT0']
      - name: PORT1
        valueFrom:
          fieldRef:
            apiVersion: v1
            fieldPath: metadata.annotations['PORT1']
      image: localhost:5000/kubetechno/consul-client
      name: kubetechno-init-consul
      volumeMounts:
      - mountPath: /mvTarget
        name: kubetechno-consul
  volumes:
    - emptyDir: {}
      name: kubetechno-consul
  containers:
  - name: server
    command:
    - ./server 
    args:
      - ${PORT0}  
    env:
    - name: kubetechno_consul_check_path
      value: /
    - name: kubetechno_consul_service
      value: kubetechno-churner-app
    - name: kubetechno_consul_buffer_seconds
      value: "10"
    - name: kubetechno_consul_timeout_seconds
      value: "30"
    - name: kubetechno_consul_period_seconds
      value: "30"
    - name: kubetechno_consul_initial_delay_seconds
      value: "3"
    - name: PORT0
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.annotations['PORT0']
    - name: PORT1
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.annotations['PORT1']
    image: localhost:5000/kubetechno/example-server
    lifecycle:
      preStop:
        exec:
          command:
          - /kubetechnoConsul/client
          - dereg
    readinessProbe:
      exec:
        command:
        - /kubetechnoConsul/client
        - check
      failureThreshold: 1
      initialDelaySeconds: 3
      periodSeconds: 30
      successThreshold: 1
      timeoutSeconds: 30
    resources:
      limits:
        kubetechno/port: "2"
      requests:
        kubetechno/port: "2"
    volumeMounts:
    - mountPath: /kubetechnoConsul
      name: kubetechno-consul
  - name: sidecar
    image: localhost:5000/kubetechno/example-sidecar
    env:
    - name: PORT0
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.annotations['PORT0']
    - name: PORT1
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.annotations['PORT1']
  hostNetwork: true
 
```

The changes that were made are:
- Annotation additions:
  - The ports assignments added to the annotations. 
  - Consul related annotations may also be added. 
  - The name of the node the pod is running on is also added.
  - A kubetechno status set to active.
- Label annotation addition of `kubetechno-node: <the pod's node>`
- An init container that copies over the consul helper binary to volume mount.
- The volume that the consul helper binary is copied to.
- The first container additions:
  - Env vars with information for consul registration and checks.
  - Env vars with port information.
  - A volume mount for the consul helper binary volume.
  - A prestop lifecycle hook to deregister.
  - A readiness probe that is tied to consul, this also registers the service instance.
  - kubetechno port resource limits and requests
- Subsequent container additions
  - Port env vars
- Setting hostNetwork to true
