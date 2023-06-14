package router

import (
	"estimation-vocabulary/internal/middlewares"
	_vocabulary "estimation-vocabulary/internal/service/vacabulary"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	//if mode == gin.ReleaseMode {
	//	gin.SetMode(gin.ReleaseMode)
	//}

	r := gin.New()
	r.Use(middlewares.Cors())
	r.Use(gin.Recovery())

	// 静态文件-App怎么部署不知道.如果也可以是dist就在这里设置
	//r.StaticFile()
	//r.GET("/api/ping", func() {
	//})
	userGroup := r.Group("/api/user")
	{
		userGroup.GET("/xxx", nil)
	}

	vocabularyGroup := r.Group("/api/vocabulary")
	{
		vocabularyGroup.GET("/xxx", _vocabulary.StartTest)
		vocabularyGroup.GET("/test", _vocabulary.Test)
		// 获得可选等级列表
		vocabularyGroup.GET("/getLevelList", _vocabulary.ShowLevelList)
		// 批处理接口-接收一个josn文件
		vocabularyGroup.POST("/batch", _vocabulary.GetScoreBatch)
	}

	return r

}
