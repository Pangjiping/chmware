package client

import (
	"fmt"
	"testing"
)

func TestLoadBalance(t *testing.T) {
	endpoints := []string{
		"aaa",
		"bbb",
		"ccc",
		"ddd",
		"eee",
		"fff",
	}

	cnts := map[string]int{}
	for i := 0; i < 1000000; i++ {
		endpoint := loadBalance(endpoints)
		cnts[endpoint]++
	}

	for k, v := range cnts {
		fmt.Println(k, float64(v)/float64(1000000))
	}
}
