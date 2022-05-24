package client

import "fmt"

// Node is node information.
// A node has many replicas by raft.
type Node struct {
	NodeID      string   `json:"node_id"`
	NodeAddress []string `json:"node_address"`
	NodeEnabled bool     `json:"node_enabled"`
}

// NodeMaps is NodeID to NodeAddress if this node is enabled.
type NodeMaps map[string][]string

var nodeMaps NodeMaps

var uniqueID map[string]bool
var uniqueIP map[string]bool

// parseNodeList parse nodeList to NodeMaps.
// It will check validation of  NodeID and NodeAddress.
func parseNodeList(nodeList []Node) error {
	for _, nodeInfo := range nodeList {
		if nodeInfo.NodeEnabled == false {
			continue
		}
		if !checkIDUnique(nodeInfo.NodeID) {
			return fmt.Errorf("invalid node ID: %s", nodeInfo.NodeID)
		}
		if nodeInfo.NodeAddress == nil || len(nodeInfo.NodeAddress) == 0 {
			return fmt.Errorf("empty address list: %s", nodeInfo.NodeID)
		}
		nodeMaps[nodeInfo.NodeID] = make([]string, 0)
		for _, addr := range nodeInfo.NodeAddress {
			// todo: 判断一个合法的IP地址 报错
			if !checkAddress(addr) {
				return fmt.Errorf("invalid IP address: %s", addr)
			}

			if !checkAddrUnique(addr) {
				return fmt.Errorf("the same IP address already exists, Address: %s", addr)
			}
			nodeMaps[nodeInfo.NodeID] = append(nodeMaps[nodeInfo.NodeID], addr)
		}
	}
	return nil
}

func init() {
	nodeMaps = make(NodeMaps)
	uniqueID = make(map[string]bool)
	uniqueIP = make(map[string]bool)
}

func checkAddress(address string) bool {
	return true
}

func checkAddrUnique(addr string) bool {
	if _, ok := uniqueIP[addr]; ok {
		return false
	}
	uniqueIP[addr] = true
	return true
}

func checkIDUnique(id string) bool {
	if _, ok := uniqueID[id]; ok {
		return false
	}
	uniqueID[id] = true
	return true
}

func tryConn(tryKey, tryVal string) error {
	return nil
}
