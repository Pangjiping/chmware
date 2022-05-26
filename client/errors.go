package client

import "errors"

var (
	ErrEmptyNodeList = errors.New("node list is empty")
	ErrPoolClosed    = errors.New("connect pool is closed")
	ErrConnClosed    = errors.New("connect is closed")
)
