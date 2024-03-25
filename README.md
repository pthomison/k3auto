# K3Auto

An golang CLI tool used for rapidly deploying kubernetes evironments in a repeatable manner for local testing

Powered By:
- [k3d](https://k3d.io)
- [flux](https://fluxcd.io/)


### Installation
```
brew tap pthomison/homebrew-tools
brew install pthomison/tools/k3auto
```

### Deployment Order Of Operations
1. Deploy k3d Cluster
2. Inject Flux Controllers
3. Deploy Embedded Deployments
4. Deploy Runtime Deployments

### Options
```
$ k3auto --help
k3auto is a local kubernetes cluster orchestrator powered by k3d and flux

Usage:
  k3auto [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create a new K3D Cluster and inject flux controllers & deployments
  delete      Delete an existing cluster
  help        Help about any command

Flags:
  -c, --cluster-config string         Override Cluster Config File
  -d, --deployment-directory string   Deployment Directory
  -h, --help                          help for k3auto
  -m, --minimal                       Only deploy the k3d cluster & flux controllers

Use "k3auto [command] --help" for more information about a command.

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


### Modifications

To embed your own deployment manifests, just fork this repository. Then add your manifests to the `default/deployments` directory and rebuild.


