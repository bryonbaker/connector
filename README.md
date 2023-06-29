# Running locally

Start the server:

`
podman run --name connector-server --rm -p 8080:8080 -e PORT=8080 -e XCHNG=2 quay.io/bryonbaker/connector-server:latest
`

Start the client. Substitute `your.host.name` for your host name. localhost will not work

`
podman run --name connector-client --rm -e NUMCONS=2 -e PORT=8080 -e URL=your.host.name -e XCHNG=2 quay.io/bryonbaker/connector-client:latest
`

# Running on OpenShift
## Setup Notes
### Number of connections
Note: To vary the load you can:

a. Play with the NUMCONS value in the client's deployment once you have it running. Specifically, change the count in this line in the `server-dep.yaml`: 

`export NUMCONS="1000"`

b. Vary the replica count in the client.

### Message Exchange Combinations
The tester can work in one of three configuraitons:
0 - Just estabish connections
1 - Clientg sends a message to server
2 - Server sends a response to the client message

The configuration is set via the `export XCHNG="0|1|2"` in the deployment config for both the server and the client apps.

## Steps
### Server Namespace
1. Deploy the server in one namespace: `oc apply -f server-dep.yaml`
2. Deploy RHSI: `skupper init --site-name server --enable-flow-collector --enable-console --console-user admin --console-password password`

### Client Namespace
1. Deploy RHSI: `skupper init --site-name client --enable-flow-collector --enable-console --console-user admin --console-password password`
2. Create the token: `skupper token create --token-type cert client-token.yaml`

### Server Namespace
1. Create the link: `skupper link create client-token.yaml`
2. Expose the deployment: `skupper expose deployment connector-server --port 8080`

### Client Namespace
1. Deploy the client in one namespace: `oc apply -f client-dep.yaml`

