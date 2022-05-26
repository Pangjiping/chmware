package client

import (
	"fmt"
	"strconv"
	"strings"
)

func init() {
	nodeMaps = make(NodeMaps)
	uniqueID = make(map[string]bool)
	uniqueIP = make(map[string]bool)
}

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
		if !nodeInfo.NodeEnabled {
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
			// 判断是否是合法的ip地址
			if !checkAddress(addr) {
				return fmt.Errorf("invalid IP address: %s", addr)
			}

			// 判断是否是唯一的ip地址
			if !checkAddrUnique(addr) {
				return fmt.Errorf("the same IP address already exists, address: %s", addr)
			}

			// 尝试通信，是否成功
			if err := tryConn(tryKey, tryVal, addr); err != nil {
				return fmt.Errorf("cannot connect server, address: %s", addr)
			}

			nodeMaps[nodeInfo.NodeID] = append(nodeMaps[nodeInfo.NodeID], addr)

		}
	}
	return nil
}

func checkAddress(address string) bool {
	// 按:拆分字符串，是否是两个
	ipAndPort := strings.Split(address, ":")
	if len(ipAndPort) != 2 {
		return false
	}

	// 按ip地址和端口号单独判断
	if !isIP(ipAndPort[0]) || !isPort(ipAndPort[1]) {
		return false
	}

	return true
}

func isIP(ip string) bool {
	if len(ip) == 0 {
		return false
	}

	if len(ip) < 7 || len(ip) > 15 {
		return false
	}

	if ip[0] == '.' || ip[len(ip)-1] == '.' {
		return false
	}

	// 按 . 分割字符串出来四个
	ss := strings.Split(ip, ".")
	if len(ss) != 4 {
		return false
	}

	for _, s := range ss {
		// 不是一个字符不能以0开头
		if len(s) > 1 && s[0] == '0' {
			return false
		}

		// 判断每个字符是否是0-9
		for i := 0; i < len(s); i++ {
			if s[i] < '0' || s[i] > '9' {
				return false
			}
		}
	}

	for i := 0; i < len(ss); i++ {
		num, err := strconv.Atoi(ss[i])
		if err != nil {
			return false
		}
		if i == 0 {
			if num < 1 || num > 255 {
				return false
			}
		} else {
			if num < 0 || num > 255 {
				return false
			}
		}
	}

	return true
}

func isPort(port string) bool {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	if portNum < 0 || portNum > 65535 {
		return false
	}

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

func tryConn(addr string, tryKey, tryVal string) error {
	// post
	req := connRequest{
		address: addr,
		key:     tryKey,
		value:   tryVal,
		method:  METHOD_POST,
	}
	postResp := invoke(req)
	if postResp.Error() != nil {
		return postResp.Error()
	}

	// get
	req = connRequest{
		address: addr,
		key:     tryKey,
		method:  METHOD_GET,
		value:   "",
	}
	getResp := invoke(req)
	if getResp.Error() != nil {
		return getResp.Error()
	}

	if getResp.Value() != tryVal {
		return fmt.Errorf("something wrong in connect trying, address: %s", addr)
	}

	// delete
	req = connRequest{
		address: addr,
		key:     tryKey,
		method:  METHOD_DELETE,
		value:   "",
	}
	delResp := invoke(req)
	if delResp.Error() != nil {
		return delResp.Error()
	}

	// get
	req.method = METHOD_GET
	reGetResp := invoke(req)
	if reGetResp.Error() == nil || reGetResp.Value() != "" {
		return fmt.Errorf("something wrong in connect trying, address: %s", addr)
	}

	return nil
}
