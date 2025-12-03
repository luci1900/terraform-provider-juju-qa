# QA

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
