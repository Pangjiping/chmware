package client

import (
	"fmt"
	"sync"

	"github.com/Pangjiping/raftdb-client/hash"
)

const (
	tryKey = "conn-test"
	tryVal = "success"
)

type Client struct {
	mu         sync.RWMutex
	consistent *hash.Consistent
	closed     bool
}

func newClient() *Client {
	return &Client{
		mu:         sync.RWMutex{},
		consistent: hash.NewConsistent(hash.DefaultNumberOfReplicas),
	}
}

func (c *Client) Set(key string, value string) error {
	if c.closed {
		return ErrConnClosed
	}
	return nil
}

func (c *Client) Get(key string) (string, error) {
	if c.closed {
		return "", ErrConnClosed
	}

	// 得到节点实例信息
	instance, err := c.consistent.Get(key)
	if err != nil {
		return "", err
	}

	// 启用负载均衡找到适合访问的ip地址
	endpoints, ok := nodeMaps[instance]
	if !ok {
		return "", fmt.Errorf("cannot found instance for key: %s", key)
	}

	addr := loadBalance(endpoints)

	resp := invoke(connRequest{
		address: addr,
		method:  METHOD_GET,
		key:     key,
	})

	return resp.Value(), nil
}

func (c *Client) Delete(key string) error {
	if c.closed {
		return ErrConnClosed
	}
	return nil
}

func (c *Client) registerInstance() error {
	if c.closed {
		return ErrConnClosed
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for nodeID, _ := range nodeMaps {
		c.consistent.Add(nodeID)
	}
	return nil
}

// Close 关闭连接资源
func (c *Client) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.closed = true
}
