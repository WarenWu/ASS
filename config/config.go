package config

var (
	Version         string
	AppName         string
	Debug           bool
	Ip              string
	Port            string
	LogLevel        string
	LogFile         string
	LogMaxAge       int //hour
	LogRotationTime int //min
	MysqlUser       string
	MysqlPassword   string
	MysqlAddr       string
	MysqlDb         string
	CrawlTimeout    int //s
	Headless        bool
)

const (
	Dubug_default           = false
	Version_defualt         = "x.x.x"
	AppName_default         = "ass_srv"
	Ip_default              = "0.0.0.0"
	Port_default            = "8080"
	ConfigFileName_defualt  = "ass_config"
	ConfigFilePath_defualt  = "./config"
	ConfigFileType_defualt  = "json"
	LogFile_default         = "./ass.log"
	LogLevel_default        = "info"
	LogRotationTime_default = 24
	LogMaxAge_default       = 30
	MysqlUser_default       = "root"
	MysqlPassword_default   = "123456"
	MysqlAddr_default       = "127.0.0.1:3306"
	MysqlDb_default         = "ass_vm"
	CrawlTimeout_default    = 30
	Headless_default        = true
)
