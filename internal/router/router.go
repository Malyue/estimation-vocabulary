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
		//建立连接，返回test_id 测试完成
		vocabularyGroup.GET("/connect", _vocabulary.StartTest)
		// 获得可选等级列表
		vocabularyGroup.GET("/getLevelList", _vocabulary.ShowLevelList)
		// 批处理接口-接收一个josn文件
		vocabularyGroup.POST("/batch", _vocabulary.GetScoreBatch)
		//获取一个当前level的单词
		vocabularyGroup.GET("/getWord", _vocabulary.GetWord)
		//用户是否认识单词处理接口
		vocabularyGroup.POST("/wordKnow", _vocabulary.UpdateLevel)
		//用户获取最终词汇量数目接口 测试完成
		vocabularyGroup.GET("/getResult", _vocabulary.GetResult)
		//测试App连通的接口 测试完成
		vocabularyGroup.GET("/exit", _vocabulary.Exit)

		vocabularyGroup.GET("/test", _vocabulary.Test)
	}

	return r

}
