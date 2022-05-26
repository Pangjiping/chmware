package client

import "math/rand"

const (
	maxRetryTimes = 3
)

// shuffle 公平洗牌算法
func shuffle(slice []int) {
	for i := len(slice); i > 0; i-- {
		lastIdx := i - 1
		idx := rand.Intn(i)
		slice[lastIdx], slice[idx] = slice[idx], slice[lastIdx]
	}
}

func loadBalance(endpoints []string) string {
	len := len(endpoints)
	indexs := make([]int, 0)
	for i := 0; i < len; i++ {
		indexs = append(indexs, i)
	}

	shuffle(indexs)
	return endpoints[indexs[0]]
}
