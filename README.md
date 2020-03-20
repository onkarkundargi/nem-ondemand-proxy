
# Nem-OnDemand-Proxy
---

## Building

To build and publish the docker image, issue the following commands:

```bash
make docker-build
make docker-push
```

You may preface the above commands with `DOCKER_REPOSITORY=` and/or `DOCKER_TAG=` to specify a custom docker repository and/or tag.

## Deployment

Nem-ondemand-proxy is deployed as a container and may be deployed using either docker-compose or kubernetes. Kubernetes is the recommended method and may be done by executing the following command:

```bash
cd kubernetes
kubectl apply -f nem-ondemand-proxy.yaml
```

To teardown the proxy, use the following:

```bash
cd kubernetes
kubectl delete -f nem-ondemand-proxy.yaml
```

In addition to deploying the container, you may also want to setup a Kubernetes port-forward. Scripts in the `kubernetes` directory are provided as an example.

## Usage

You will need to know the device id of the ONU you wish to query. This may be done with `voltctl device list`.

There is no dedicated client yet, and it is recommended to use `grpcurl`. For example,

```bash
# perform on-demand query on device with id 2910b26bbb29521d93fab21b
grpcurl -plaintext -d '{"id": "2910b26bbb29521d93fab21b"}' localhost:50052 on_demand_api.NemService/OmciTest
```
