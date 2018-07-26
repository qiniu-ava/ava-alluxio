package utils

import (
	"os"
)

func GetWalkerName() (string, error) {
	return os.Hostname()
}
