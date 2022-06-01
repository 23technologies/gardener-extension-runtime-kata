# Gardener Extension for Katacontainers

[![reuse compliant](https://reuse.software/badge/reuse-compliant.svg)](https://reuse.software/)

[Project Gardener](https://github.com/gardener/gardener) implements the automated management and operation of [Kubernetes](https://kubernetes.io/) clusters as a service.
Its main principle is to leverage Kubernetes concepts for all of its tasks.

This controller implements Gardener's extension contract for the `kata` ContainerRuntime.
The latest release's `ControllerRegistration` resource that can be used to register this controller to Gardener can be found [here](https://github.com/23technologies/gardener-extension-runtime-kata/releases/latest/download/controller-registration.yaml).

Please find more information regarding the extensibility concepts and a detailed proposal [here](https://github.com/gardener/gardener/blob/master/docs/proposals/01-extensibility.md).


<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Gardener Extension for Katacontainers](#gardener-extension-for-katacontainers)
- [What does this package provide?](#what-does-this-package-provide)
- [Wait! How does this differ from kata-deploy?](#wait-how-does-this-differ-from-kata-deploy)
- [Current limitations](#current-limitations)
- [How to...](#how-to)
	- [Enable the extension as a gardener operator](#enable-the-extension-as-a-gardener-operator)
	- [Use it as end user](#use-it-as-end-user)

<!-- markdown-toc end -->

# What does this package provide?
Generally, this extension should enable using [katacontainers](https://katacontainers.io/) as container runtime within `Shoot` clusters.
Thus, you can isolate your workload easily into [qemu](https://www.qemu.org/) / [firecracker](https://firecracker-microvm.github.io/) / [cloudhypervisor](https://www.cloudhypervisor.org/) VMs on the nodes, when this extension is installed.

# Wait! How does this differ from kata-deploy?
Actually, it doesn't.
Internally, the extension is using the same mechanism as provided by [kata-deploy](https://github.com/kata-containers/kata-containers/tree/main/tools/packaging/kata-deploy).
However, when it comes to using katacontainers in combination with firecracker VMs, the machine image needs to be tweaked a bit.
See also this [how-to](https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/how-to-use-kata-containers-with-firecracker.md) to get an idea of what needs to be configured for the usage of katacontainers with firecracker.
This configuration overhead is handled by this extension.
For this, the extension controller modifies the `OperatingSystemConfiguration` via a `mutatingWebhookConfiguration` so that the preparation of the nodes is performed.
More precisely, an additional script for the devmapper setup is provided and executed by an additional systemd service unit.
Moreover, the controller deploys some helm-charts containing the configuration and the daemonset provided by kata-deploy.

# Current limitations
Generally, the limitations are two-fold.
First, the cloud service provider your `Shoot` is running on needs to support nested virtualization for katacontainers to work properly.
On the Azure cloud, nested virtualization is enabled for various [machine types](https://docs.microsoft.com/en-us/azure/virtual-machines/acu).
As our development and testing environment was limited to `Shoots` hosted on Azure, the extension controller will only do something meaningful, when `Shoots` are spawned on Azure.
Second, the implementation was only tested for Ubuntu as machine image.
Thus, the controller will idle, when other machine images are selected.
Of course, the support for other machine images can be implemented.
However, this is a future topic, and this package should be considered as proof of concept.

# How to...

## Enable the extension as a gardener operator
Of course, you need to apply the `controller-registration` resources to the garden cluster first.
You can find the corresponding yaml-files in our [releases](https://github.com/23technologies/gardener-extension-runtime-kata/releases).
Moreover, you will need to modify your Azure-cloudprofile so that it contains the following configuration:
``` yaml
...
  machineImages:
  - name: ubuntu
	versions:
	- cri:
	  - containerRuntimes:
		- type: kata
		name: containerd
	  - name: docker
	  version: 18.4.20210415
...
```

## Use it as end user
If you have access rights to the garden cluster, go ahead and follow the instructions in [Enable the extension as a gardener operator](#enable-the-extension-as-a-gardener-operator).
Otherwise, ask your gardener operator to do so.
Once the extension is installed, you should be able to select `kata` as additional runtime from the Gardener dashboard during shoot creation or simply specify the corresponding container runtime in your `Shoot.Spec`.
To run some example workload in your `Shoot` cluster you can use the  [sample workloads](https://github.com/kata-containers/kata-containers/blob/main/tools/packaging/kata-deploy/README.md#run-a-sample-workload) from kata-deploy.
