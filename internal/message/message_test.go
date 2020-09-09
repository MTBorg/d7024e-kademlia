package message_test

import (
	"errors"
	"fmt"
	"kademlia/internal/kademliaid"
	"kademlia/internal/message"
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
	msg := message.New(content, target)

	assert.Equal(t, msg.Target, target)
	assert.Equal(t, msg.Content, content)
}

func TestDeserialize(t *testing.T) {
	var msg message.Message
	var err error
	// Should return an empty message and error if the string is empty
	msg, err = message.Deserialize("")
	assert.Empty(t, msg)
	assert.Error(t, err)

	// Should return an empty message and error if the string only contains a
	// sender id and no separator
	msg, err = message.Deserialize("senderid")
	assert.Empty(t, msg)
	assert.Error(t, err)

	// Should be able to pass empty content
	msg, err = message.Deserialize("senderid;")
	assert.NoError(t, err)
	assert.Equal(t, msg.Content, "")
}

func TestSend(t *testing.T) {
	testId := strings.Repeat("1", 40) //IDs are 160-bit (= 40 hex characters)
	var senderMock *SenderMock
	msg := message.Message{SenderId: kademliaid.FromString(testId), Content: "content", Target: "target"}
	msgSerialized := fmt.Sprintf("%s;content", testId)
	var err error

	// Should return the error from send if there was an error
	senderMock = new(SenderMock)
	senderMock.On("Send", msgSerialized).Return(errors.New("this is an error"))
	err = msg.Send(senderMock)
	assert.Equal(t, err, errors.New("this is an error"))
	senderMock.AssertExpectations(t)

	// Should return nil if send does not return an error
	senderMock = new(SenderMock)
	senderMock.On("Send", msgSerialized).Return(nil)
	err = msg.Send(senderMock)
	assert.NoError(t, err)
	senderMock.AssertExpectations(t)
}
