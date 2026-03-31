# QA
 
## Locally

Bootstrap two new controllers:
```shell
juju bootstrap microk8s tfqa
juju bootstrap microk8s tfqa-offering
```

Run all tests:
```shell
make test
```

Or run only a specific test:
```shell
make run=PrivateRegistry test
```

## SolQA cluster

Go to https://github.com/canonical/terragrunt-deployment-pipelines/actions/workflows/maas_physical.yaml
Click "Run workflow".
Pick a cluster, ideally one that doesn't already have a workflow currently running.
Pick "solution", since we're bootstrapping our own controllers.
Pass in parameters including the repo and branch you want to run, like:
    {"repo": "canonical/terraform-provider-juju-qa", "ref": "main"}
Click "Run workflow"

# Tags

Every resource (Juju controller, application, etc.) must have tags attached to ensure it runs against the correct MAAS cluster.

Tags look like `category,cluster`, as an example: `juju_upgrade,sqa-dh1_j8_1`.

An inventory tag is available, including both VMs and metal. Two VM tags are prepared for running Juju controllers, `juju` and `juju_upgrade`.
