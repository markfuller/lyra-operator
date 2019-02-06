# Lyra operator example

Experiments with [operator-sdk](https://github.com/operator-framework/operator-sdk). You'll need to install the 'operator-sdk' command line tool to use this repo.

## Pre-requistes

Unlike most Lyra repos, this project doesn't (and won't) support Go modules so you will need to disable support:

```bash
    export GO111MODULE=off
```

## Usage

Before first time use, you need to install the Workflow CRD:

```bash
    kubectl create -f deploy/crds/lyra_v1alpha1_workflow_crd.yaml
```

Create a Workflow resource like this:

```bash
    kubectl apply -f deploy/crds/lyra_v1alpha1_workflow_sample.yaml
```

The workflow resource can be deleted by either of:

```bash
    kubectl delete -f deploy/crds/lyra_v1alpha1_workflow_sample.yaml
    kubectl delete workflows sample2-workflow
```

## Development

You can run the operator (with mock applicator event hooks) directly from this repo like this:

```bash
    operator-sdk up local --namespace=default
```

If you make changes you might need to regenerate the controller code:

```bash
    operator-sdk generate k8s
```
