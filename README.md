# kube-remote-provisioner

*Work in progress* Remote storage provisioner for Kubernetes.

## Idea

The idea is to implement a Kubernetes storage provisioner that manages volumes on a remote fileserver. Plan is to have two parts:

* An agent running on a fileserver, managing volumes and setting up NFS exports
* A storage provider based on [nfs-provisioner](https://github.com/kubernetes-incubator/nfs-provisioner) which talks to agents over gRPC, and mounting the volumes over NFS.

The first plans are to implement an agent running on ZFS capable fileservers.
