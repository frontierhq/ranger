.. image:: https://github.com/gofrontier-com/ranger/actions/workflows/ci.yml/badge.svg
    :target: https://github.com/gofrontier-com/ranger/actions/workflows/ci.yml

|

.. image:: logo.png
  :width: 200
  :alt: Ranger logo
  :align: center

======
Ranger
======

Ranger is a command line tool that enables "GitOps for everything" by providing a technology-agnostic contract-based framework to define how workloads (and their associated configuration and secrets) are built, tested, deployed (and destroyed) and promoted through environments.

.. contents:: Table of Contents
    :local:

-----
About
-----

Ranger has been built to allow platform teams to provide a deployment framework

Ranger enables common GitOps outcomes - including improved efficiency and security, a better developer experience, reduced costs, and faster deployments - all from your existing CI/CD platform.

--------
Concepts
--------

Video?

Mermaid diagram?

~~~~~~
GitOps
~~~~~~

Something about GitOps and it's traditional use in K8S apps, but how Ranger enables "for everything".

~~~~~~~~~
Workloads
~~~~~~~~~

Something about reusable units of business value (?) that have a common interface regardless of implementation, i.e. build, test, deploy, check, destroy.

~~~~
Sets
~~~~

Something about what a set represents, like application, user journey, tier/stack?

--------
Download
--------

~~~~~~~
Release
~~~~~~~

Binaries and packages of the latest stable release are available at `https://github.com/gofrontier-com/ranger/releases <https://github.com/gofrontier-com/ranger/releases>`_.

~~~~~~~~~
Extension
~~~~~~~~~

The Ranger extension for Azure DevOps is available from `Visual Studio Marketplace <https://marketplace.visualstudio.com/items?itemName=gofrontier.ranger>`_, which will automatically install Ranger via a task.

-------------
Configuration
-------------

~~~~~~~~
Examples
~~~~~~~~

Example workload and set manifest?

.. code:: yaml

  ---
  version: 23
  environment: dev
  set: core-infra
  nextEnvironment: prd
  workloads:
    - name: microsoft-defender
      type: Shared/microsoft-defender-workload
      version: 1.5.1
    - name: virtual-network
      type: Shared/virtual-network-workload
      verson: 3.1.7
    - name: sql-server
      type: Shared/sql-server-workload
      version: 2.9.3
    - name: app-gateway
      type: Shared/app-gateway-workload
      version: 1.1.8

.. code:: yaml

  ---
  version: 9
  environment: dev
  set: creditcards-infra
  nextEnvironment: prd
  workloads:
    - name: kubernetes-cluster
      type: Shared/kubernetes-cluster-workload
      version: 6.0.3
    - name: app-gateway-ingress
      type: Shared/app-gateway-ingress-workload
      version: 1.0.2
    - name: api-gateway-service
      type: Shared/api-gateway-service-workload
      version: 11.3.1
    - name: statements-service
      type: CreditCardsLZ/statements-service-workload
      version: 1.7.3


.. code:: yaml

  ---
  version: 17
  environment: dev
  set: currentaccounts-infra
  nextEnvironment: prd
  workloads:
    - name: kubernetes-cluster
      type: Shared/kubernetes-cluster-workload
      version: 6.0.3
    - name: app-gateway-ingress
      type: Shared/app-gateway-ingress-workload
      version: 1.0.2
    - name: api-gateway-service
      type: Shared/api-gateway-service-workload
      version: 11.3.1
    - name: withdrawal-service
      type: CurrentAccountsLZ/withdrawal-service-workload
      version: 5.1.9

------------
Contributing
------------

We welcome contributions to this repository. Please see `CONTRIBUTING.md <https://github.com/gofrontier-com/ranger/tree/main/CONTRIBUTING.md>`_ for more information.
