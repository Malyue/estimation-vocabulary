package main

import (
	_config "estimation-vocabulary/config"
	_router "estimation-vocabulary/internal/router"
)

func main() {
	//_algorithm.ImportToDb()
	config := _config.GetConfig()
	//_model.InitMysql(config)
	//_algorithm.ImportToDb()
	r := _router.Init()
	r.Run(config.Server.Addr)
}
