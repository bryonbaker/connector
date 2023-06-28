# Running locally

Start the server:

`
podman run --name connector-server -p 8080:8080 --rm  connector-server:latest
`

Start the client. Substitute `your.host.name` for your host name. localhost will not work

`
podman run --name connector-client --rm -e NUMCONS=100 -e URL=your.host.name connector-client:latest
`