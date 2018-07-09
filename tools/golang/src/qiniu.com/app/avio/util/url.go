package util

import (
	"net/url"
	"path"
	"regexp"
	"strings"
)

var zoneRP = regexp.MustCompile("//(z[0-9]){0, 1}/")

func Url2LocalPath(fpath, prefix string) string {
	if strings.HasPrefix(fpath, "http") {
		u, e := url.Parse(fpath)
		if e != nil {
			return ""
		}
		return path.Join(prefix, u.Path)
	} else if strings.HasPrefix(fpath, "qiniu") {
		zoneRP.ReplaceAllString(fpath, "//")
		u, e := url.Parse(fpath)
		if e != nil {
			return ""
		}

		return path.Join(prefix, u.Path)
	}
	return path.Join(prefix, fpath)
}
