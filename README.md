# d7024e-kademlia

### Running

#### Compose

There are two variants of the services, dev (development) and prod (production).
The dev service supports hot reloading (whereas prod does not) but has a larger
image size than the prod variant.
What variant is used can be specified by passing either
`docker-compose-dev.yml` or `docker-compose-prod.yml` to the relevant
`<compose-file>` option.

To run a single node using compose run

```bash
docker-compose -f <compose-file> up
```

#### Swarm

Make sure you have a cluster initialized by running

```bash
docker swarm init
```

Then run

```bash
docker stack deploy --compose-file <compose-file> kademlia
```

### Executing commands

Any node can told to execute any command by running the following

```
docker exec -ti <container-name/id> kademliactl <command>
```

### Testing

Tests can be run by running

```bash
go test ./test/...
```
