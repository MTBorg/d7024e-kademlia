package httplistener_test

import (
	"kademlia/internal/datastore"
	httplistener "kademlia/internal/http/listener"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

type HttpWriterMock struct {
	mock.Mock
}

func (m *HttpWriterMock) Write(data []byte) (int, error) {
	args := m.Called(data)
	return args.Int(0), nil
}

func (mock *HttpWriterMock) WriteHeader(statusCode int) {
}

func (mock *HttpWriterMock) Header() http.Header {
	return nil
}

func TestHandleGet(t *testing.T) {
	var writerMock *HttpWriterMock
	var n node.Node
	var handler httplistener.RequestHandler
	var req *http.Request

	//Should tell user about missing hash
	writerMock = new(HttpWriterMock)
	req, _ = http.NewRequest("GET", "/objects/", nil)

	n = node.Node{}
	handler = httplistener.RequestHandler{Node: &n}
	writerMock.On("Write", []byte("Missing hash")).Return(0, nil)
	handler.HandleGet(writerMock, req)
	writerMock.AssertExpectations(t)

	// Should return the value if found
	n = node.Node{}
	writerMock = new(HttpWriterMock)
	writerMock.On("Write", mock.Anything).Return(0, nil)
	n.DataStore = datastore.New()
	msg := "this is a message"
	hash := kademliaid.NewKademliaID(&msg)
	req, _ = http.NewRequest("POST", "/objects/"+hash.String(), nil)
	n.DataStore.Insert(msg)
	handler = httplistener.RequestHandler{Node: &n}
	handler.HandleGet(writerMock, req)
	writerMock.AssertCalled(t, "Write", []byte(msg+", from local node"))
}
