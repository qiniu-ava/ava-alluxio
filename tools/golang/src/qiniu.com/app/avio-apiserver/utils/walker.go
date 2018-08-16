package utils

import (
	"qiniu.com/app/common/typo"
)

func PickWalker(walkers []typo.Walker) string {
	// TODO consider use more complex algorithm to pick walker
	// for there maybe jobs cost every long time
	if len(walkers) == 0 {
		// we should never be here
		return ""
	}
	maxIndex := 0
	max := 0
	for i, w := range walkers {
		if len(w.Jobs) > max {
			maxIndex = i
			max = len(w.Jobs)
		}
	}

	return string(walkers[maxIndex].Name)
}
