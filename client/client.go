package client

import (
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
	return "", nil
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
