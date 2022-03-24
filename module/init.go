package module

import (
	"ASS/config"
	"ASS/model"
	"ASS/router"
)

func InitModule() {
	config.InitConfig()
	router.InitHttpRouter()
	model.InitDatabse()
}
