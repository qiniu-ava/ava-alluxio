package main

import (
	"qiniu.com/app/common/typo"
	log "qiniupkg.com/x/log.v7"
)

type Executor interface {
	Do(path string) error
}

type PreloadExecutor struct {
	client *AlluxioClient
}

func NewPreloadExecutor(client *AlluxioClient) (p *PreloadExecutor) {
	return &PreloadExecutor{client}
}

func (p *PreloadExecutor) Do(path string) error {
	log.Infof("get kafka message to preload file from %s", path)
	return p.client.Call(typo.PreloadJobType, path)
}

type SaveExecutor struct {
	client *AlluxioClient
}

func NewSaveExecutor(client *AlluxioClient) (p *SaveExecutor) {
	return &SaveExecutor{client}
}

func (p *SaveExecutor) Do(path string) error {
	log.Infof("get kafka message to save file from %s", path)
	return p.client.Call(typo.SaveJobType, path)
}
