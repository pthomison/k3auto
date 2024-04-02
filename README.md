# K3Auto

An golang CLI tool used for rapidly deploying kubernetes evironments in a repeatable manner for local testing

Powered By:
- [flux](https://fluxcd.io/)
- [k3d](https://k3d.io)
- [k3s](https://k3s.io/)
- [k8s](https://kubernetes.io/)


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
- [metrics-server](https://github.com/kubernetes-sigs/metrics-server)
- [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)
- [capacitor](https://github.com/gimlet-io/capacitor)


To *only* deploy a the k3d cluster and flux controllers, use the `--minimal`/`-m` flag.


The default k3d cluster config file is embeded in the binary, but can be found at `k3auto/default/k3d-config.yaml`. To use a different config file at runtime, use the `--cluster-config`/`-c` flag.

### Deploying Resources

To deploy your own desired resources at runtime, use the `--deployment-directory`/`-d` flag and supply yaml manifests within that directory. At present moment, k3auto will capture that directory into an OCI Image, ship the image to the k3d registry, then create flux OCIRepository & Kustomization objects that will deploy your manifests into your cluster. A kustomization.yaml can be supplied within this directory if any kustomization changes are desired.


### Modifications

To embed your own deployment manifests, just fork this repository. Then add your manifests to the `default/deployments` directory and rebuild.


### Roadmap
1. An `update` subcommand for refreshing user deployments
2. Better image solution for OCI deployments (network registry layer is a current weak point)
3. Standardize ingress solution