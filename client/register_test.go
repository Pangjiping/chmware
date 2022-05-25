package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getNodeList() []Node {
	return []Node{
		{
			NodeID:      "node-0",
			NodeAddress: []string{"127.0.0.1:3000", "127.0.0.1:3001", "127.0.0.1:3002"},
			NodeEnabled: true,
		},
		{
			NodeID:      "node-1",
			NodeAddress: []string{"127.0.0.1:3003", "127.0.0.1:3004", "127.0.0.1:3005"},
			NodeEnabled: true,
		},
		{
			NodeID:      "node-2",
			NodeAddress: []string{"127.0.0.1:3006", "127.0.0.1:3007", "127.0.0.1:3008"},
			NodeEnabled: true,
		},
		{
			NodeID:      "node-4",
			NodeAddress: []string{"127.0.0.1:3009", "127.0.0.1:3010", "127.0.0.1:3011"},
			NodeEnabled: false,
		},
	}
}

func TestParseNodeList(t *testing.T) {
	nodeList := getNodeList()

	// no error
	err := parseNodeList(nodeList)
	require.NoError(t, err)

	// repeated ip address
	repeatedAddress := []Node{
		{
			NodeID:      "node-5",
			NodeAddress: []string{"127.0.0.1:3005"},
			NodeEnabled: true,
		},
	}
	err = parseNodeList(repeatedAddress)
	require.Error(t, err, "repeated node ip address")

	// repeated node id
	repeatedNodeID := []Node{
		{
			NodeID:      "node-2",
			NodeAddress: []string{},
			NodeEnabled: true,
		},
	}
	err = parseNodeList(repeatedNodeID)
	require.Error(t, err, "repeated node id")

	// empty node address list
	emptyAddressList := []Node{
		{
			NodeID:      "node-empty",
			NodeAddress: []string{},
			NodeEnabled: true,
		},
	}
	err = parseNodeList(emptyAddressList)
	require.Error(t, err, "empty node address list")

	// invalid ip address
	// todo
}

func TestIsIP(t *testing.T) {
	validIP := "127.0.0.1"
	invalidIPOne := "1.1.1"
	invalidIPTwo := "0.1.1.1"
	invalidIPThree := "9.0.0.566"
	invalidIPFour := "a.b.c.d"

	valid := isIP(validIP)
	require.Equal(t, true, valid)

	valid = isIP(invalidIPOne)
	require.Equal(t, false, valid)

	valid = isIP(invalidIPTwo)
	require.Equal(t, false, valid)

	valid = isIP(invalidIPThree)
	require.Equal(t, false, valid)

	valid = isIP(invalidIPFour)
	require.Equal(t, false, valid)
}

func TestCheckAddress(t *testing.T) {
	validAddress := "127.0.0.1:3306"
	invalidAddressOne := "127.0.0.1:99999"
	invalidAddressTwo := "127.0.0.1:"
	invalidAddressThree := "127.0.0.256:11"

	valid := checkAddress(validAddress)
	require.Equal(t, true, valid)

	valid = checkAddress(invalidAddressOne)
	require.Equal(t, false, valid)

	valid = checkAddress(invalidAddressTwo)
	require.Equal(t, false, valid)

	valid = checkAddress(invalidAddressThree)
	require.Equal(t, false, valid)
}
