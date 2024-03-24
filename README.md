# K3Auto

An golang CLI tool used for rapidly deploying kubernetes evironments in a repeatable manner for local testing

Powered By:
- k3d
- flux

### Installation
```
brew tap pthomison/homebrew-tools
brew install pthomison/tools/k3auto
```

### Options
```
$ k3auto --help
k3auto is a local kubernetes cluster orchestrator powered by k3d and flux

Usage:
  k3auto [flags]

Flags:
  -c, --cluster-config string         Override Cluster Config File
  -d, --deployment-directory string   Deployment Directory
  -h, --help                          help for k3auto
  -m, --minimal                       Only deploy the k3d cluster & flux controllers
```

###  Cluster Deployment

By default, k3auto will deploy a single node k3d cluster and will inject the following resources
```
metrics-server
```

To *only* deploy a the k3d cluster and flux controllers, use the `--minimal`/`-m` flag.

The default k3d cluster config file is embeded in the binary, but can be found at `k3auto/default/k3d-config.yaml`. To use a different config file at runtime, use the `--cluster-config`/`-c` flag.

### Deploying Resources

To deploy your own desired resources at runtime, use the `--deployment-directory`/`-d` flag and supply yaml manifests within that directory. At present moment, k3auto will attempt to cast the objects to know types and then use the ctrl client to apply it into the cluster. There is a known limitation with deploying additional CRDs that is actively being worked on.


### Options & Modifications



