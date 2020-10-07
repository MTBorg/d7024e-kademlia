package httplistener

import (
	"fmt"
	"io/ioutil"
	cmdparser "kademlia/internal/command/parser"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"path"

	"github.com/rs/zerolog/log"

	"net/http"
)

// Wrap the handler in a struct so that we can pass the node object to it.
// This so stupid, but so is go.
type RequestHandler struct {
	Node *node.Node
}

func (handler *RequestHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	log.Trace().Msg("Received GET request")
	hash := path.Base(r.URL.String())

	// If hash is equal to object then no hash was passed in
	if hash == "objects" {
		log.Info().Msg("Received GET request missing hash")
		w.Write([]byte("Missing hash"))
		return
	}

	cmd := cmdparser.ParseCmd(fmt.Sprintf("get %s", hash))
	value, err := cmd.Execute(handler.Node)
	if err != nil {
		log.Error().Msg("Failed to execute get command")
	}
	w.Write([]byte(value))
}

func (handler *RequestHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	log.Trace().Msg("Received POST request")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Msgf("Error reading request body %s", err.Error())
		return
	}

	cmd := cmdparser.ParseCmd(fmt.Sprintf("put %s", string(data)))
	value, err := cmd.Execute(handler.Node)
	if err != nil {
		log.Error().Str("Error", err.Error()).Msg("Failed to execute put command")
		w.Write([]byte(err.Error()))
		return
	}

	// Add location header and set status code to 201 (created)
	sData := string(data)
	hash := kademliaid.NewKademliaID(&sData)
	w.Header().Add("Location", "/objects/"+hash.String())
	w.WriteHeader(201)

	w.Write([]byte(value))
}

func Listen(node *node.Node) error {
	requesthandler := RequestHandler{Node: node}

	// Register handlers
	http.HandleFunc("/objects/", requesthandler.HandleGet)
	http.HandleFunc("/objects", requesthandler.handlePost)

	// ListenAndServe always returns non-nil error and blocks until it does
	err := http.ListenAndServe(":8080", nil)
	return err
}
