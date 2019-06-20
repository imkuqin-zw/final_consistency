package consumer

import (
	cluster "github.com/bsm/sarama-cluster"
	"sync"
)

var (
	s        *cluster.Consumer
	initOnce sync.Once
	nftOnce  sync.Once
	errOnce  sync.Once
)

func init() {
	initOnce.Do(func() {
		s = newClusterConsumer()
	})
}

func GetConsumer() *cluster.Consumer {
	return s
}

func SetNtfHandle(handle ntfHandle) {
	nftOnce.Do(func() {
		go notificationProcess(handle)
	})
}

func SetErrHandle(handle errHandle) {
	errOnce.Do(func() {
		go errProcess(handle)
	})
}
