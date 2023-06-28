# Running locally

Start the server:

`
podman run --name connector-server -p 8080:8080 --rm  connector-server:latest
`

Start the client. Substitute `your.host.name` for your host name. localhost will not work

`
podman run --name connector-client --rm -e NUMCONS=100 -e URL=your.host.name connector-client:latest
`

# Running on OpenShift
Note: Play with the NUMCONS value in the client's deployment once you have it running.

## Server Namespace
1. Deploy the server in one namespace: `oc apply -f server-dep.yaml`
2. Deploy RHSI: `skupper init --site-name server --enable-flow-collector --enable-console --console-user admin --console-password password`

## Client Namespace
1. Deploy RHSI: `skupper init --site-name client --enable-flow-collector --enable-console --console-user admin --console-password password`
2. Create the token: `skupper token create --token-type cert client-token.yaml`

## Server Namespace
1. Create the link: `skupper link create client-token.yaml`
2. Expose the deployment: `skupper expose deployment connector-server --port 8080`

## Client Namespace
1. Deploy the client in one namespace: `oc apply -f client-dep.yaml`

