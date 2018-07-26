package common

import (
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func ZKElect(servers []string, path string) (string, error) {
	conn, _, err := zk.Connect(servers, time.Second*60, nil)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	// TODO retry
	b, _, err := conn.Get(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
