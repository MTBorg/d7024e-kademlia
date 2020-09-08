package rpc_test

import (
	"errors"
	"fmt"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SenderMock struct {
	mock.Mock
}

func (m *SenderMock) Send(data string) error {
	args := m.Called(data)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	var content, target = "some message", "127.0.0.1:1337"
	rpc := rpc.New(content, target)

	assert.Equal(t, rpc.Target, target)
	assert.Equal(t, rpc.Content, content)
}

func TestDeserialize(t *testing.T) {
	var r rpc.RPC
	var err error
	// Should return an empty message and error if the string is empty
	r, err = rpc.Deserialize("")
	assert.Empty(t, r)
	assert.Error(t, err)

	// Should return an empty message and error if the string only contains a
	// sender id and no separator
	r, err = rpc.Deserialize("senderid")
	assert.Empty(t, r)
	assert.Error(t, err)

	// Should be able to pass empty content
	r, err = rpc.Deserialize("senderid;rpcid;")
	assert.NoError(t, err)
	assert.Equal(t, r.Content, "")
}

func TestSend(t *testing.T) {
	testId := strings.Repeat("1", 40) //IDs are 160-bit (= 40 hex characters)
	var senderMock *SenderMock
	rpc := rpc.RPC{SenderId: kademliaid.FromString(testId), RPCId: kademliaid.FromString(testId), Content: "content", Target: "target"}
	rpcSerialized := fmt.Sprintf("%s;%s;content", testId, testId)
	var err error

	// Should return the error from send if there was an error
	senderMock = new(SenderMock)
	senderMock.On("Send", rpcSerialized).Return(errors.New("this is an error"))
	err = rpc.Send(senderMock)
	assert.Equal(t, err, errors.New("this is an error"))
	senderMock.AssertExpectations(t)

	// Should return nil if send does not return an error
	senderMock = new(SenderMock)
	senderMock.On("Send", rpcSerialized).Return(nil)
	err = rpc.Send(senderMock)
	assert.NoError(t, err)
	senderMock.AssertExpectations(t)
}
