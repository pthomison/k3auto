# K3Auto

An golang CLI tool used for rapidly deploying kubernetes environments in a repeatable manner for local testing

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
3. Inject Docker Registry
4. Inject Secrets
5. Deploy Embedded Deployments
6. Deploy Runtime Deployments

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
  update      Reinject deployments
  version     Prints the version, commit, & build date

Flags:
  -c, --cluster-config string         Override Cluster Config File
  -d, --deployment-directory string   Deployment Directory
  -h, --help                          help for k3auto
  -m, --minimal                       Only deploy the k3d cluster & flux controllers
  -s, --secret-config string          Inject Secrets To the Cluster on Creation

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

### Updating Resources

After creating your environment, you can update the manifests inside your deployment directory and then run `k3auto update`. This will update the manifests inside the environment and then deploy your changes.


### Secrets

With local kubernetes clusters, the need inevitably arises to inject & manage secrets for passwords, cloud credentials, etc. Storing passwords in plaintext in a repo is a non-starter, so you could inject them by hand after the local environment is created but this quickly becomes cumbersome.

Generally I prefer to use a secret store external to my environment (such as AWS Parameter Store), then use automation to inject that into my environment (like the wonderful [External Secrets](https://external-secrets.io/) project). This automation still requires bootstrapping however, so I still need a way to inject secrets from my local machine at runtime.

To support this, k3auto uses a secrets config file that will support different secret backends and will automatically create these as secrets in your cluster at creation.

Currently, the only secret backend that is written is a bash exec backend (this is to enable shelling out to other CLI tools such as the awscli or accessing the systems keyring). If tigheter integration with a given provider is required, a custom backend can be written trivialy, with the only interface function being `Resolve(ctx context.Context, args []string) (string, error)`.

Example Secrets Config File
```
DefaultSecret: defaultSecretName
DefaultNamespace: default
Secrets:
  - Type: exec
    SecretName: "flux-system"
    SecretKey: "known-hosts"
    Args:
      - /bin/bash
      - -c
      - "aws ssm get-parameter --name /lab/flux-known-hosts --with-decryption | jq -r '.Parameter.Value'"

  - Type: exec
    SecretName: "flux-system" # supports adding multiple keys to a single secret
    SecretKey: "flux-public-key"
    Args:
      - /bin/bash
      - -c
      - "aws ssm get-parameter --name /lab/flux-public-key --with-decryption | jq -r '.Parameter.Value'"

```


### Modifications

To embed your own deployment manifests, just fork this repository. Then add your manifests to the `default/deployments` directory and rebuild.


### Roadmap
1. ~~An `update` subcommand for refreshing user deployments~~ (Update now works!)
2. ~~Better image solution for OCI deployments (network registry layer is a current weak point)~~ (Now registry lives within the cluster && k3auto port-forwards into the pod for connectivity)
3. Standardize ingress solution
4. Preloading base images to prevent excessive bandwidth consumption
5. Secrets Management (In Progress)
6. Internal package testing
7. Flux Reconciliation Requests


### Project Learnings
- Port-Forwarding from CLI apps is a little complication/not well documented, but is extremely useful for k8s environments in unknown network setups
- Go makes glueing together a lot of different projects/tools ~~very~~ relatively easy, all dependencies between the projects will need to be compatible. Not too big of an issue for unrelated tools, but can cause some funkiness if they have similar domains (k3d and docker-cli being the combo I ran into on this project)
- Garbage Collected Languages (ie golang) need/should/(good idea?) to invoke a non-GC lang for cryptographic key operations so it can release the key memory and not rely on GC to remove it. Still working on groking this one, but checkout https://github.com/containers/image/issues/1634 if interested
- Executing `kubectl apply` like commands from client-go is possible (should be... since kubectl is written with it :wink: ), but the command does a lot of work under the covers & I have found it easier to just decode the objects (a la `runtime.Decoder`) and then use the standard client-go semantics.
- When splitting yaml objects (specifically for CRDs) on `---`, be sure to use regex/start-of-line/end-of-line