package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery()) // self-recovery
	router.Use(gin.Logger())   // Custom logging middleware
	router.Use(cors.Cors()) // Custom cors middleware

	v1G := router.Group("/v1")
	{
		addCommon(v1G)
	}

	return router
}

// func addCommon(router *gin.RouterGroup) {
	/////////////   Common   /////////////
	iCommon := common.NewCommonAPI()
	commonG := router.Group("common")
	{
		commonG.POST("/images/upload", iCommon.ImagesUpload)
		commonG.POST("/resources/upload", iCommon.ResourcesUpload)
	}
}
