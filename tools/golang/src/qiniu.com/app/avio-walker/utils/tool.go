package utils

import (
	"fmt"
	"os"

	"qiniu.com/app/avio-walker/compile"
)

func GetWalkerHost(port int) (string, error) {
	hn, e := os.Hostname()
	if e != nil {
		return "", nil
	}

	host := hn + ".avio-walker-svc.ava.svc.cluster.local"
	if compile.Mode == "debug" {
		host = "0.0.0.0"
	}
	if port != 0 {
		return fmt.Sprintf("%s:%d", host, port), e
	}
	return fmt.Sprintf("%s", host), e
}
