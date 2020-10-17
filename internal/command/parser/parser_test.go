package cmdparser_test

import (
	"kademlia/internal/command"
	"kademlia/internal/command/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCmd(t *testing.T) {
	var cmd command.Command

	//should return nil if the parsing failed
	cmd = cmdparser.ParseCmd("invalid command")
	assert.Nil(t, cmd)

	//should be able to parse a ping command
	// TODO: Should also test that target is set
	cmd = cmdparser.ParseCmd("ping localhost")
	assert.NotNil(t, cmd)

	//should be able to parse a message command
	// TODO: Should also test that target and content is set
	cmd = cmdparser.ParseCmd("get somehash")

	// should be able to parse a exit command
	cmd = cmdparser.ParseCmd("exit")
	assert.NotNil(t, cmd)

	// should be able to parse a storage command
	cmd = cmdparser.ParseCmd("storage")
	assert.NotNil(t, cmd)

	// should be able to parse a getid command
	cmd = cmdparser.ParseCmd("getid")
	assert.NotNil(t, cmd)

	// should be able to parse a addcontact command
	// TODO: Should also test that nodeID and address is set
	cmd = cmdparser.ParseCmd("addcontact nodeid address")
	assert.NotNil(t, cmd)

	// should be able to parse a getcontacts command
	cmd = cmdparser.ParseCmd("getcontacts")
	assert.NotNil(t, cmd)

	// should be able to parse a init command
	// TODO: Should also test that address is set
	cmd = cmdparser.ParseCmd("init address")
	assert.NotNil(t, cmd)

	// should be able to parse a put command
	cmd = cmdparser.ParseCmd("put")
	assert.Nil(t, cmd)

	// should be able to parse a join command
	cmd = cmdparser.ParseCmd("join")
	assert.NotNil(t, cmd)

	//should return nil if an invalid command was passed
	cmd = cmdparser.ParseCmd("non-existent command")
	assert.Nil(t, cmd)

	//should return nil if invalid options were passed
	cmd = cmdparser.ParseCmd("ping") //ping needs a target option
	assert.Nil(t, cmd)
}
