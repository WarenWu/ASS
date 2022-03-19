package module

import (
	"ASS/config"
	"ASS/db"
	"ASS/router"
)

func InitModule() {
	config.InitConfig()
	router.InitHttpRouter()
	db.InitDatabse()
}
