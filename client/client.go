package client

import (
	"sync"

	"github.com/Pangjiping/chmware/hash"
)

type Client struct {
	mu         sync.RWMutex
	consistent hash.Consistent
}
