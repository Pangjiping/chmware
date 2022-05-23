package client

import (
	"sync"

	"github.com/Pangjiping/raftdb-client/hash"
)

type Client struct {
	mu         sync.RWMutex
	consistent hash.Consistent
}
