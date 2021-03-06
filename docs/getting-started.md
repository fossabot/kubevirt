# Getting Started

A quick start guide to get KubeVirt up and running inside our container based
development cluster.

## I just want it built and run it on my cluster

First, point the `Makefile` to the docker registry of your choice:

```bash
export DOCKER_PREFIX=index.docker.io/myrepo
export DOCKER_TAG=mybuild
```

Then build the manifests and images:

```bash
make && make push
```

Finally push the manifests to your cluster:

```bash
kubectl create -f _out/manifests/release/kubevirt-operator.yaml
kubectl create -f _out/manifests/release/kubevirt-cr.yaml
```

### Docker Desktop for Mac
The bazel build system doesn't support the macOS keychain. Please ensure that
you deactivate the option `Securely store Docker longins in macOS keychain` in
the Docker preferences. After restarting the docker service login with `docker
login`. Your `$HOME/.docker/config.json` should look like:

```json
{
  "auths" : {
    "https://index.docker.io/v1/" : {
      "auth" : "XXXXXXXXXX"
    }
  },
  "credSstore" : ""
}
```

## Building

The KubeVirt build system runs completely inside docker. In order to build
KubeVirt you need to have `docker` and `rsync` installed. You also need to have `docker`
running, and have the [permissions](https://docs.docker.com/install/linux/linux-postinstall/#manage-docker-as-a-non-root-user) to access it.

**Note:** For running KubeVirt in the dockerized cluster, **nested
virtualization** must be enabled - [see here for instructions for Fedora](https://docs.fedoraproject.org/en-US/quick-docs/using-nested-virtualization-in-kvm/index.html).
As an alternative [software emulation](software-emulation.md) can be allowed.
Enabling nested virtualization should be preferred.

### Dockerized environment

Runs master and nodes containers, when each one of them run virtual machine via QEMU.
In additional it runs dnsmasq and docker registry containers.

### Compatibility

The minimum compatible Kubernetes version is 1.11.0. Important features required
for scheduling and memory are missing or incompatible with previous versions.

### Compile and run it

To build all required artifacts and launch the
dockerizied environment, clone the KubeVirt repository, `cd` into it, and:

```bash
# Build and deploy KubeVirt on Kubernetes in our vms inside containers
export KUBEVIRT_PROVIDER=k8s-1.13.3 # this is also the default if no KUBEVIRT_PROVIDER is set
make cluster-up
make cluster-sync
```

This will create a virtual machine called `node01` which acts as node and master. To create
more nodes which will register themselves on master, you can use the
`KUBEVIRT_NUM_NODES` environment variable. This would create a master and one
node:

```bash
export KUBEVIRT_NUM_NODES=2 # schedulable master + one additional node
make cluster-up
```

You can use the `KUBEVIRT_MEMORY_SIZE` environment 
variable to increase memory size per node. Normally you don't need it, 
because default node memory size is set.

```bash
export KUBEVIRT_MEMORY_SIZE=8192M # node has 8GB memory size
make cluster-up
```

To destroy the created cluster, type

```
make cluster-down
```

**Note:** Whenever you type `make cluster-down && make cluster-up`, you will
have a completely fresh cluster to play with.

### Accessing the containerized nodes via ssh

Based on the used cluster, node names might be different.
You can get the names from following command:

```bash
# cluster-up/kubectl.sh get nodes
NAME                        STATUS   ROLES    AGE   VERSION
kind-1.14.2-control-plane   Ready    master   11h   v1.14.2
```

Then you can execute the following command to access the node:
```
# ./cluster-up/ssh.sh kind-1.14.2-control-plane
root@kind-1:/#
```

### Automatic Code Generation

Some of the code in our source tree is auto-generated (see `git ls-files|grep '^pkg/.*generated.*go$'`).
On certain occasions (but not when building git-cloned code), you would need to regenerate it
with

```bash
make generate
```

Typical cases where code regeneration should be triggered are:

 * When changing APIs, REST paths or their comments (gets reflected in the api documentation, clients, generated cloners...)
 * When changing mocked interfaces (the mock generator needs to update the generated mocks then)

 We have a check in our CI system, so that you don't miss when `make generate` needs to be called.

### Testing

After a successful build you can run the *unit tests*:

```bash
    make
    make test
```

They do not need a running KubeVirt environment to succeed.
To run the *functional tests*, make sure you have set
up a dockerized environment. Then run

```bash
    make cluster-sync # synchronize with your code, if necessary
    make functest # run the functional tests against the dockerized VMs
```

If you'd like to run specific functional tests only, you can leverage `ginkgo`
command line options as follows (run a specified suite):

```
    FUNC_TEST_ARGS='-ginkgo.focus=vmi_networking_test -ginkgo.regexScansFilePath' make functest
```

In addition, if you want to run a specific test or tests you can prepend any `Describe`,
`Context` and `It` statements of your test with an `F` and Ginkgo will only run items
that are marked with the prefix. Remember to remove the prefix before issuing
your pull request.

For additional information check out the [Ginkgo focused specs documentation](http://onsi.github.io/ginkgo/#focused-specs)

## Use

Congratulations you are still with us and you have built KubeVirt.

Now it's time to get hands on and give it a try.

### Create a first Virtual Machine

Finally start a VMI called `vmi-ephemeral`:

```bash
    # This can be done from your GIT repo, no need to log into a VMI

    # Create a VMI
    ./cluster-up/kubectl.sh create -f examples/vmi-ephemeral.yaml

    # Sure? Let's list all created VMIs
    ./cluster-up/kubectl.sh get vmis

    # Enough, let's get rid of it
    ./cluster-up/kubectl.sh delete -f examples/vmi-ephemeral.yaml


    # You can actually use kubelet.sh to introspect the cluster in general
    ./cluster-up/kubectl.sh get pods

    # To check the running kubevirt services you need to introspect the `kubevirt` namespace:
    ./cluster-up/kubectl.sh -n kubevirt get pods
```

This will start a VMI on master or one of the running nodes with a macvtap and a
tap networking device attached.

#### Example

```bash
$ ./cluster-up/kubectl.sh create -f examples/vmi-ephemeral.yaml
vm "vmi-ephemeral" created

$ ./cluster-up/kubectl.sh get pods
NAME                              READY     STATUS    RESTARTS   AGE
virt-launcher-vmi-ephemeral9q7es  1/1       Running   0          10s

$ ./cluster-up/kubectl.sh get vmis
NAME           LABELS                        DATA
vmi-ephemera    kubevirt.io/nodeName=node01   {"apiVersion":"kubevirt.io/v1alpha2","kind":"VMI","...

$ ./cluster-up/kubectl.sh get vmis -o json
{
    "kind": "List",
    "apiVersion": "v1",
    "metadata": {},
    "items": [
        {
            "apiVersion": "kubevirt.io/v1alpha2",
            "kind": "VirtualMachine",
            "metadata": {
                "creationTimestamp": "2016-12-09T17:54:52Z",
                "labels": {
                    "kubevirt.io/nodeName": "master"
                },
                "name": "vmi-ephemeral",
                "namespace": "default",
                "resourceVersion": "102534",
                "selfLink": "/apis/kubevirt.io/v1alpha2/namespaces/default/virtualmachineinstances/testvm",
                "uid": "7e89280a-be62-11e6-a69f-525400efd09f"
            },
            "spec": {
    ...
```

### Accessing the Domain via VNC

First make sure you have `remote-viewer` installed. On Fedora run

```bash
dnf install virt-viewer
```

Windows users can [download remote-viewer from virt-manager.org](https://virt-manager.org/download/), and may need
to add virt-viewer installation folder to their `PATH`.

Then, after you made sure that the VMI `vmi-ephemeral` is running, type

```
cluster-up/virtctl.sh vnc vmi-ephemeral
```

to start a remote session with `remote-viewer`.

`cluster-up/virtctl.sh` is a wrapper around `virtctl`. `virtctl` brings all
virtual machine specific commands with it and is a supplement to `kubectl`.

**Note:** If accessing your cluster through ssh, be sure to forward your X11 session in order to launch `virtctl vnc`.

### Bazel and KubeVirt

#### Build.bazel merge conflicts

You may encounter merge conflicts in `BUILD.bazel` files when creating pull
requests. Normally you can resolve these conflicts extremely easy by simply
accepting the new upstream version of the files and run `make` again. That will
update the build files with your changes.

