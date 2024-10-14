Migration Cluster Day2

This repository is meant to provide to means for Day2 operations on a target cluster for VM migrations.

Some or all the items here may end up as an integral part of MTV, and for the time being we will make all efforts
to automate and smooth the experience to the maximum we can.


# Using this repo to prepare a migration cluster

## Prerequisite

- Installed cluster
  - 3 node baremetals that can run virtualization
  - each node has extra disk for storage
- Operators
  - OpenShift Virtualization
  - LVM Storage (TODO try to automate that installation with the helm chart)
- `helm` client installed
- `oc` client installed

## Install

- clone this repo
- install the helm chart

```console
helm install migration chart/migartion-cluster-day2 
```

## Initialize the MTV provider

Navigate to the mtv-init application route
```console
oc get route -n default mtv-init -o jsonpath={.status.ingress[0].host}
```

Fill in the details of form:
- vddk image: go and download the vddk.tar.gz from broadcom site. The link is part of the form
- vcenter username: use the admin username, or a user which is the most credentials you can get
- vcenter password
- vcenter url: the url of vcenter, no need to add /sdk in the end

When the form is submitted, follow the job in the default namespace that creates the vddk image and updates the existing 
vmware-credentials


Check the vmware provider to see if it can connect to the vsphere instance.
```console
oc get provider -n openshift-mtv
```

