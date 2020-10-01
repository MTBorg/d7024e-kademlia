CONT_NAME="^kademlia_kademlia." # Common part of container names

echo "Initilizing all running Kademlia containers..."

# Get all container IDs
cont_ids=$(docker ps -aq -f name=$CONT_NAME -f status=running) 

# Init each node with their assigned IPs
for id in $cont_ids; do
  cont_ip="$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $id)"

  docker exec -it $id kademliactl init $cont_ip > /dev/null
done

# Insert a known node in the network to each routing table
known_node_cid="$(echo "$cont_ids" | head -n 1)"
known_node_ip="$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $known_node_cid)"
known_node_kadid="$(docker exec -it $known_node_cid kademliactl getid | sed -ne 's/^.*response: //p')" 
for id in $cont_ids; do
  if [ "$id" != "$known_node_cid" ]; then
    docker exec -it $id sh -c "\
			kademliactl addcontact "$known_node_kadid" "$known_node_ip" && \
			kademliactl join" > /dev/null
  fi
done

echo "Done"
