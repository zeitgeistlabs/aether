# Compute Aether


Compute Aether (pronounced like "ether", lasting name still TBD) aims to solve common storage and coordination problems for distributed services,
by building transparent abstractions into the platform.

The Aether API eventually aims to expose basic utilities, such as:

* Queues
* Key-value stores
* Locks
* Election

Aether itself handles the resources and logic required. Administrators do not need to provision resources for individual services.

**Project status: Lab**

(This project is experimental in nature, will regularly experience breaking changes, and is not production-quality.)

# Getting Started

## Cluster Setup

Enable [service account token projection](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#service-account-token-volume-projection) on your cluster.

## Compute Aether Setup

Deploy Compute Aether with `kubectl create -f deploy.yaml`. Review the service account permissions to make sure they are acceptable.

## Service Setup

1. Create a unique service account for your service.
2. Add a projected token to your service spec.

# Building


There is a crude Makefile for the project.

`make binary` builds a local binary (fast). Requires gomod enabled.

`make image` builds a hermetic Docker image (slow).

## Testing


## Contributing

Please file issues with suggestions! While you're welcome to make PRs, at this point the project is in flux,
and it will be hard to contribute code without coordination.

References
==

[1] https://jpweber.io/blog/a-look-at-tokenrequest-api/
