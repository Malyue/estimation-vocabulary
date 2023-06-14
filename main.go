package main

import (
	_config "estimation-vocabulary/config"
	_internal "estimation-vocabulary/internal"
	_router "estimation-vocabulary/internal/router"
	"sync"
)

func main() {
	config := _config.GetConfig()
	r := _router.Init()
	// 初始化全局map
	_internal.UserMap = sync.Map{}
	// 开启定时任务定期清除没有被删除的map
	//go _pkg.DeleteMap(_internal.UserMap)
	r.Run(config.Server.Addr)
}
