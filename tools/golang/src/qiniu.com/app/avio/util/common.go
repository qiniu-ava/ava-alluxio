package util

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"qiniu.com/app/common/typo"
)

func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(b)/float64(div), "KMGTPE"[exp])
}

func DurationString(d time.Duration) string {
	if d.Hours() > 48 {
		return fmt.Sprintf("%.0f days", math.Ceil(d.Hours()/24))
	} else if d.Hours() > 0 {
		return fmt.Sprintf("%.0f hours", math.Ceil(d.Hours()))
	} else if d.Minutes() > 0 {
		return fmt.Sprintf("%.0f mins", math.Ceil(d.Minutes()))
	} else {
		return fmt.Sprintf("%.0f seconds", math.Ceil(d.Seconds()))
	}
}

func validatePath(p string) error {
	if strings.HasPrefix(p, "/workspace/mnt/bucket/") {
		return nil
	}
	return fmt.Errorf("【%s】不是 bucket 中的数据", p)
}

func ResolveAlluxioPath(localPath string, uid int64) (string, error) {
	if strings.HasPrefix(localPath, "/") {
		if e := validatePath(localPath); e != nil {
			return "", e
		}
		return strings.Replace(localPath, "/workspace/mnt/bucket/", fmt.Sprintf("/ava/qn-bucket/%d/", uid), 1), nil
	}

	d, e := os.Getwd()
	if e != nil {
		return "", fmt.Errorf("解析本地地址出错: %v", e)
	}
	u := path.Join(d, localPath)
	if e = validatePath(u); e != nil {
		return "", e
	}
	return strings.Replace(u, "/workspace/mnt/bucket/", fmt.Sprintf("/ava/qn-bucket/%d/", uid), 1), nil
}

func NewJobName(jobType typo.JobType) string {
	rand.Seed(time.Now().Unix())
	f := rand.Float32()
	return fmt.Sprintf("%s-%x", jobType.String(), math.Float32bits(f))
}
