![Go](https://github.com/cloudish-ufcg/k8s-qos-driven-plugins/workflows/Go/badge.svg)

# k8s-qos-driven-plugins
This repository hosts the code required to use and deploy a kubernetes cluster using a scheduler based on QoS-driven scheduling policy.

# Deployment
We are using the version v1.18.3 of **Kubernetes**, so some features may or may not work in different versions.

Deploying a simple kubernetes cluster can be done [using minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/#starting-a-cluster).
To create a local kubernetes cluster using the version **v1.18.3** and virtualization driver **kvm** simply run: 

```bash
minikube start --kubernetes-version v1.18.3 --driver=kvm2 --bootstrapper=kubeadm
```

To activate a plugin already registered your [scheduler image](Kube-Scheduler) you need to provide [KubeSchedulerConfiguration](KubeSchedulerConfiguration) to it.
The controllers using the defined profile(s) will use this scheduler.

In the cluster directory you can see how to deploy a kubernetes cluster and use it to execute the experiment scenarios.

## KubeSchedulerConfiguration
This should configure a *profile* which allows **QosAware** plugin to run.
Create one (or just adapt [ours](scheduler-config.conf)) and place it at `/etc/kubernetes/scheduler-config.conf` in the master node - it's where our [deploy](kube-scheduler.yaml) expects to read from. 

The *profile* have to enable our plugin in **queueSort**, **preBind** and **postBind** sections:  
```yaml
plugins:
  queueSort:
    enabled:
      - name: QosAware
    disabled:
      - name: "*"
  preBind:
    enabled:
      - name: QosAware
  postBind:
    enabled:
      - name: QosAware
```

> **_NOTE:_** One and only one plugin queueSort plugin should be enabled for each profile, this is why we explicitly disable the default one 

## Kube Scheduler
**Kube-scheduler** is actually the default scheduler of kubernetes. This project (mostly) compiles into a version of **kube-scheduler** with our plugin registered.  

The **kubelet** has a parameter `--pod-manifest-path` that usually points to `/etc/kubernetes/manifests` path on the master node.
There should have a `kube-scheduler.yaml` file. The **kubelet** reads this file to deploy the cluster's scheduler.

To deploy a scheduler using our image create a deployment - like [this one](kube-scheduler.yaml) - and place into the master node at `/etc/kubernetes/manifests/kube-scheduler.yaml`.
The kubelet will detect the changes and deploy the new scheduler.

# Extending kubernetes
The kubernetes has behaviours that can be customized with plugins or extensions.
We change specific behaviours of the [scheduling framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/) with plugins.

## Install extra plugins
In [deploy section](KubeSchedulerConfiguration) is shown how to configure plugins, but we can only activate plugins that are registered.
The kube-scheduler already has default plugins registered, but to register a new one we have to generate a version of the scheduler that binds it.
This binding is exactly what happens in [our main file](cmd/main.go).   

## Creating a new plugin
To your plugin implement any extension point, the `Plugin` object that [is binded](Install-extra-plugins) have to implement their respective interface.
The interfaces can [be found here](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/v1alpha1/interface.go) and they look like this:

```go
type Plugin interface {
    Name() string
}

type QueueSortPlugin interface {
    Plugin
    Less(*v1.pod, *v1.pod) bool
}

type PreFilterPlugin interface {
    Plugin
    PreFilter(context.Context, *framework.CycleState, *v1.pod) error
}

// ...
```

## Building a new image
You can build a new kube-scheduler image by building a program that runs the scheduler with the plugins.go).

Running `CGO_ENABLED=0 go build -o bin/kube-scheduler cmd/main.go` will create a executable of [this program](cmd/main.go) and put it in `bin/kube-scheduler`.
You can then put this executable into your docker image and use it to deploy.  

# Creating roles to run customized scheduler
The plugins developed in this project need to monitor some properties of the current active controllers (deployments and jobs) in the cluster. By default, the user `system:kube-scheduler` does not have permission to get information about `job` resources. In order to allow the user `system:kube-scheduler` to get the required information simply run:

```
kubectl create -f cluster/deploy/reader_job/reader-job-role.yaml 
kubectl create -f cluster/deploy/reader_job/reader-job-rolebindings.yaml
```

# Debugging
In the code are some logs messages.
Their visibility can be controlled by changing the `--v` argument on [scheduler deployment](kube-scheduler.yaml).

There's also a **debug tool you can access at http://localhost:10000/**:
```console
foo@bar:~$ kubectl --namespace=kube-system exec kube-scheduler-minikube -- curl -s http://localhost:10000/
          CONTROLLER    QOS/SLO                   TTV               RUNNING               WAITING               BINDING
                      1.00/0.80                 31m1s                2h4m6s                    0s                    0s
            50_slo_1  0.44/0.50            -2m22.541s                 8m13s                10m36s                 591ms
            50_slo_2  0.48/0.50              -34.946s                 9m22s                 9m57s                 882ms
```

This tool (such as some **kubectl** commands) works well with `watch` command.
So one can run `watch -n 2 -- kubectl --namespace=kube-system exec kube-scheduler-minikube -- curl -s http://localhost:10000/` to monitoring controllers' QoS.
