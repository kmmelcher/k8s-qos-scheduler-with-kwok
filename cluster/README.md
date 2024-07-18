# The Kubernetes cluster deployment

In order to deploy a Kubernetes cluster, you just need to follow the below steps:

## Step 1: Set up environment (all nodes)

You need to install the Docker and the Kuberntes tools in all nodes of the cluster. We have two scripts that help you in this step. You should execute the following commands:

- `sudo bash deploy/docker-install.sh`
- `sudo bash deploy/kubernetes-install.sh`

## Step 2: Master deployment (master node)

In case of using bare-instances to deploy the cluster, it is necessary to disable the swap memory on each server. If you are using virtualized instances to deploy the Kubernetes cluster, you do not need to execute the following command.

- `sudo swapoff -a`

Start the master service by executing the following command.

- `sudo kubeadm init --pod-network-cidr=10.244.0.0/16`

Once this command finishes, it will display a **kubeadm join** message at the end. This will be used to join the worker nodes to the cluster.

Next, enter the following to create a directory for the cluster. It will allow the current user to use the Kubernetes client.
- `mkdir -p $HOME/.kube`
- `sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config`
- `sudo chown $(id -u):$(id -g) $HOME/.kube/config`

A Pod Network is a way to allow communication between different nodes in the cluster. In our cluster, we use the flannel virtual network.

Enter the following:

- `sudo kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml`

Verify that everything is running and communicating:

- `kubectl get pods --all-namespaces`

## Step 3: Worker deployment (worker nodes)

Again, in case of using bare-instances to deploy the cluster, it is necessary to disable the swap memory on each server. If you are using virtualized instances to deploy the Kubernetes cluster, you do not need to execute the following command.

- `sudo swapoff -a`

Now you need to enter the **kubeadm join** command on each worker node to connect it to the cluster. If you did not make note of this while you were starting the master component, you can generate the **kubeadm join** entering the following on master node:

- `sudo kubeadm token create --print-join-command`

You can check if it works by checking the status of the nodes. Switch to the master server, and enter:

- `kubectl get nodes`

The system should display the worker nodes that you joined to the cluster.

# Creating priority classes in the Kubernetes cluster

In order to execute the experiment scenarios using the default scheduling policy, our experiment design consider some classes of priorities. For this reason, we need to create these classes of priorities before to run the experiments. To do this, tun the following commands on master node:

- `kubectl create -f conf/priority_class/prod-priority-class.yaml`
- `kubectl create -f conf/priority_class/batch-priority-class.yaml`
- `kubectl create -f conf/priority_class/free-priority-class.yaml`
- `kubectl create -f conf/priority_class/opportunistic-priority-class.yaml`

# Creating roles to run customized scheduler

The plugins developed in this project (to be used in customized QoS-driven scheduling scnearios) need to monitor some properties of the current active controllers (deployments and jobs) in the cluster. By default, the user `system:kube-scheduler` does not have permission to get information about `job` resources. In order to allow the user `system:kube-scheduler` to get the required information simply run:

```
kubectl create -f deploy/reader_job/reader-job-role.yaml 
kubectl create -f deploy/reader_job/reader-job-rolebindings.yaml
```

# Running the broker to submit a workload

## Install project dependencies

We are running the broker application on master node. It is necessary to install the dependencies of broker. First of all, tou need to install the go as follows:


```
   wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
   sudo tar -xvzf go1.14.4.linux-amd64.tar.gz
   sudo mv go /usr/local
   
   export GOROOT=/usr/local/go
   export PATH=$GOROOT/bin:$PATH
```

The above example is installing the `1.14.4` version of go. 
Additionally, the above environment variables will be set for your current session only. To make it permanent add above commands in `~/.profile` file.

In order to download the project dependencies, run the following command inside the project root directory.

   ```
   go mod download
   ```

## Running the broker

TODO

Inside of broker directory you can run:

```
   bash run_experiment1.sh
   ```

## Broker output

TODO