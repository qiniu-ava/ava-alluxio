package main

import (
	"qiniu.com/app/avio-apiserver/services"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	jobService := services.NewJobService()

	router.GET("/jobs/:name", jobService.GetJobInfo)
	router.GET("/jobs", jobService.ListJobs)
	router.POST("/jobs", jobService.CreateJob)
	router.POST("/jobs/:name/start", jobService.StartJob)
	router.PUT("/jobs/:name", jobService.UpdateJobInfo)
	router.DELETE("/jobs/:name", jobService.DeleteJobInfo)

	return router
}
