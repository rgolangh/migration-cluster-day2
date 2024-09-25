# migration-cluster-preprep

- custom manifest
A manifest that will be applied as part of the assisted installer setup
this manifest will install the helm chart and potentially some more resources which are do not
fit in the helm chart


- helm chart
installs the following assets:
  - local storage
    - make sure storage class for kubevirt is annotated

  - mtv operator
      -  configure mtv host
      -  configure provider (depends on vsphere credentials)
        check if provider has condition for connection to vsphere. if not, add it. 

  - volume with vddk tar.gz
  - build config for vddk
  - vddk image
 
  - network
    - install nmstate operator

# 
