package config

import (
	"io"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Option struct {
	ConfigFileName string
	ConfigFilePath string
	ConfigFileType string
}

var option Option

func init() {
	pflag.BoolVar(&Debug, "debug", false, "exec pattern")
	pflag.StringVar(&Port, "port", Port_default, "listen port")
	pflag.StringVar(&Ip, "ip", Ip_default, "listen ip")
	pflag.StringVar(&option.ConfigFileName, "configFileName", ConfigFileName_defualt, "configfile name")
	pflag.StringVar(&option.ConfigFilePath, "configFilePath", ConfigFilePath_defualt, "configfile path")
	pflag.StringVar(&option.ConfigFileType, "configFileType", ConfigFileType_defualt, "configfile fmt")
}

func InitConfig() {
	setViper()
	readConfig()
	InitLog()
}

func setDefault() {
	viper.SetDefault("log.file", LogFile_default)
	viper.SetDefault("log.level", LogLevel_default)
	viper.SetDefault("log.rotationtime", LogRotationTime_default)
	viper.SetDefault("log.maxage", LogMaxAge_default)
	viper.SetDefault("mysql.addr", MysqlAddr_default)
	viper.SetDefault("mysql.password", MysqlPassword_default)
	viper.SetDefault("mysql.user", MysqlUser_default)
	viper.SetDefault("mysql.db", MysqlDb_default)
	viper.SetDefault("version", Version_defualt)
	viper.SetDefault("appname", AppName_default)
	viper.SetDefault("crawltimeout", CrawlTimeout_default)
	viper.SetDefault("headless", Headless_default)
	viper.SetDefault("cn_stock_pe_max", MaxStockPe_default)
	viper.SetDefault("cn_stock_pe_min", MinStockPe_default)
	viper.SetDefault("cn_pe_max", MaxPe_default)
	viper.SetDefault("cn_pe_min", MinPe_default)
}

func setViper() {
	setDefault()
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	viper.SetConfigName(option.ConfigFileName)
	viper.SetConfigType(option.ConfigFileType)
	viper.AddConfigPath(option.ConfigFilePath)
	err := viper.ReadInConfig() //根据上面配置加载文件
	if err != nil {
		logrus.Warnf("Don't find config file !\n")
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		logrus.Infoln(in.String())
		readConfig()
		InitLog()
	})
}

func readConfig() {
	Debug = viper.GetBool("debug")
	Version = viper.GetString("version")
	AppName = viper.GetString("appname")
	Ip = viper.GetString("ip")
	Port = viper.GetString("port")
	MysqlAddr = viper.GetString("mysql.addr")
	MysqlPassword = viper.GetString("mysql.password")
	MysqlUser = viper.GetString("mysql.user")
	MysqlDb = viper.GetString("mysql.db")
	CrawlTimeout = viper.GetInt("crawltimeout")
	Headless = viper.GetBool("headless")
	MaxStockPe = viper.GetInt("cn_stock_pe_max")
	MinStockPe = viper.GetInt("cn_stock_pe_min")
	MaxPe = viper.GetInt("cn_pe_max")
	MinPe = viper.GetInt("cn_pe_min")
}

func InitLog() {
	LogFile = viper.GetString("log.file")
	LogLevel = viper.GetString("log.level")
	LogRotationTime = viper.GetInt("log.rotationtime")
	LogMaxAge = viper.GetInt("log.maxage")
	if Debug {
		logrus.SetReportCaller(true)
	}
	logrus.SetLevel(Level(LogLevel))
	logrus.SetFormatter(&logrus.JSONFormatter{})
	fileWriter, err := rotatelogs.New(
		LogFile+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(LogFile),
		rotatelogs.WithMaxAge(time.Duration(LogMaxAge)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(LogRotationTime)*time.Minute),
	)
	if err != nil {
		logrus.Errorf("failed to create rotatelogs: %s\n", err)
	} else {
		stdoutWriter := os.Stdout
		logrus.SetOutput(io.MultiWriter(stdoutWriter, fileWriter))
	}
}

func Level(l string) (level logrus.Level) {
	switch l {
	case "trace":
		level = logrus.TraceLevel
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	}
	return
}
