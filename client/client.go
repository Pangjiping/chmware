package client

import (
	"log"
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
}

func NewClient(numberOfReplicas int) *Client {
	if numberOfReplicas == 0 {
		numberOfReplicas = hash.DefaultNumberOfReplicas
	}
	return &Client{
		mu:         sync.RWMutex{},
		consistent: hash.NewConsistent(numberOfReplicas),
	}
}

func (c *Client) Set(key string, value string) error {
	return nil
}

func (c *Client) Get(key string) (string, error) {
	return "", nil
}

func (c *Client) Delete(key string) error {
	return nil
}

func (c *Client) RegisterInstance(nodeList []Node) error {
	if nodeList == nil || len(nodeList) == 0 {
		return ErrEmptyNodeList
	}

	if len(nodeList) == 1 {
		log.Printf("Only one kv instance is applied: %v", nodeList[0])
	}

	if err := parseNodeList(nodeList); err != nil {
		return err
	}
	return nil
}
