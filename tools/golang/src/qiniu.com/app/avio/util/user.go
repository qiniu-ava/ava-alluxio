package util

import (
	"path"

	"qiniupkg.com/x/config.v7"
)

const ConfigDir string = "/workspace/config"

type UserInfo struct {
	UID          int64    `json:"uid"`
	Name         string   `json:"name,omitempty"`
	Bucket       string   `json:"bucket,omitempty"`
	Key          string   `json:"key,omitempty"`
	Secret       string   `json:"secret,omitempty"`
	PublicKeys   []string `json:"publicKeys,omitempty"`
	Privileged   bool     `json:"privileged,omitempty"`
	Groups       []string `json:"groups,omitempty"`
	DefaultGroup string   `json:"defaultGroup,omitempty"`
	SpecVersion  string   `json:"specVersion,omitempty"`
}

var userConf UserInfo

func GetUserConfig() (info UserInfo, e error) {
	info = UserInfo{}
	e = config.LoadFile(&info, path.Join(ConfigDir, "user"))
	return
}

func GetUID() (uid int64, e error) {
	if info, e := GetUserConfig(); e == nil {
		return info.UID, nil
	} else {
		return 0, e
	}
}
