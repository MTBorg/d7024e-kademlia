OUTPUT_PATH=bin/
go build -o $OUTPUT_PATH/kademlia cmd/kademlia/kademlia.go && \
	go build -o $OUTPUT_PATH/kademliactl cmd/kademliactl/kademliactl.go
