# Omni Infrastructure Provider for SUSE's [Harvester](https://github.com/harvester/harvester)

Can be used to automatically provision Talos nodes in a harvester cluster.

## Running Infrastructure Provider

First you need to create a service account for the infrastructure provider.

```bash
$ omnictl serviceaccount create --role=InfraProvider --use-user-role=false infra-provider:harvester

Set the following environment variables to use the service account:
OMNI_ENDPOINT=https://<account-name>.omni.siderolabs.io/
OMNI_SERVICE_ACCOUNT_KEY=<service-account-key>

Note: Store the service account key securely, it will not be displayed again
```

Create a service account kubeconfig for your harvester cluster.
Store it in `kubeconfig` file.

### Using Docker

```bash
docker run -it -d -v ./kubeconfig:/kubeconfig ghcr.io/siderolabs/omni-infra-provider-harvester --kubeconfig-file /kubeconfig --omni-api-endpoint https://<account-name>.omni.siderolabs.io/ --omni-service-account-key <service-account-key>
```

### Using Executable

Build the project (should have docker and buildx installed):

```bash
make omni-infra-provider-linux-amd64
```

Run the executable:

```bash
_out/omni-infra-provider-linux-amd64 --kubeconfig-file kubeconfig --omni-api-endpoint https://<account-name>.omni.siderolabs.io/ --omni-service-account-key <service-account-key>
```
