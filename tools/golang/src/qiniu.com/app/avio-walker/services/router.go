package services

import (
	"github.com/gin-gonic/gin"
	"qiniu.com/app/avio-walker/walker"
)

func InitRouter(walker walker.Walker) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	walkJobService := NewWalkJobService(walker)

	router.POST("/jobs/:name/start", walkJobService.StartWalkJob)

	return router
}
