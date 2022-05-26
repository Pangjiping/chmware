package client

import (
	"fmt"
	"log"
	"sync"
)

type ConnPool struct {
	mu       sync.Mutex
	connChan chan *Client
	factory  func() *Client
	closed   bool
	nodeList []Node
}

// RegisterConn 注册一个连接池，包括了新建连接池和向每个连接中注册实例
func RegisterConn(nodeList []Node, size uint) (*ConnPool, error) {
	var (
		connPool *ConnPool
		err      error
	)
	if connPool, err = newPool(newClient, size); err != nil {
		return nil, err
	}

	connPool.mu.Lock()
	defer connPool.mu.Unlock()

	connPool.nodeList = nodeList

	// registe nodeList
	if err := connPool.registeNodeList(); err != nil {
		return nil, err
	}

	for conn := range connPool.connChan {
		if err = conn.registerInstance(); err != nil {
			connPool.Close()
			return nil, err
		}
	}

	return connPool, nil
}

// newPool 新建连接池
func newPool(fn func() *Client, size uint) (*ConnPool, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid size: %d", size)
	}

	return &ConnPool{
		mu:       sync.Mutex{},
		factory:  fn,
		connChan: make(chan *Client, size),
	}, nil
}

// GetConn 获取连接
func (p *ConnPool) GetConn() (*Client, error) {
	select {
	case r, ok := <-p.connChan:
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		client := p.factory()
		err := client.registerInstance()
		if err != nil {
			return nil, err
		}
		return client, nil
	}
}

// Close 关闭连接池
func (p *ConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	// 关闭连接chan
	close(p.connChan)

	// 关闭通道里的资源
	for c := range p.connChan {
		c.close()
	}
}

// Release 释放某个连接
func (p *ConnPool) Release(c *Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 如果连接池关闭了，直接释放
	if p.closed {
		c.close()
		return
	}

	select {
	case p.connChan <- c:
	default:
		c.close()
	}
}

func (p *ConnPool) registeNodeList() error {
	if p.closed {
		return ErrPoolClosed
	}

	if p.nodeList == nil || len(p.nodeList) == 0 {
		return ErrEmptyNodeList
	}

	if len(p.nodeList) == 1 {
		log.Printf("only one kv instance is applied: %v", p.nodeList[0])
	}

	if err := parseNodeList(p.nodeList); err != nil {
		return err
	}
	return nil
}
