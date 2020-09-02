# d7024e-kademlia

### Running

#### Compose

To run a single node using compose run

```bash
docker-compose up
```

#### Swarm

Make sure you have a cluster initialized by running

```bash
docker swarm init
```

Then run

```bash
docker stack deploy --compose-file docker-compose.yml kademlia
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
