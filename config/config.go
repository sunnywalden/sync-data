package config

import (
	"flag"
	"github.com/sirupsen/logrus"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf global
	Conf *TomlConfig
	//LogLevel logrus.Level
)

type TomlConfig struct {
	App AppInfo `toml:"app"`
	OA OA `toml:"oa"`
	Redis RedisConf `toml:"redis"`
	Mysql MysqlConf `toml:"mysql"`
	Log  LogConf `toml:"log"`
}

type AppInfo struct {
	Name string `toml:"name"`
	Url string `toml:"url"`
	Port string `toml:"port"`
	Debug bool `toml:"debug"`
	Version string `toml:"version"`
}

type OA struct {
	OaUrl string `toml:"oa_url"`
	TokenApi string `toml:"token_api"`
	UserApi string `toml:"user_api"`
	OaClient string `toml:"oa_client"`
	OaSecret string `toml:"oa_secret"`
	UserActive uint8 `toml:"user_active"`
}

type RedisConf struct {
	Host string `toml:"cache_host"`
	Port string `toml:"cache_port"`
	DB int `toml:"cache_database"`
	Password string `toml:"cache_password"`
}

type MysqlConf struct {
	Host string `toml:"db_host"`
	Port string `toml:"db_port"`
	DB string `toml:"db_database"`
	User string `toml:"db_user"`
	Password string `toml:"db_password"`
}

type LogConf struct {
	Level logrus.Level `toml:"log_level"`
}

func init() {
	log.Printf("Configure init!\n")
	flag.StringVar(&confPath, "conf", "./config/application.toml", "-conf path")
}

// Init, init conf
func Init() (*TomlConfig, error) {
	_, err := toml.DecodeFile(confPath, &Conf)
	if err != nil {
		log.Fatalf("getting configure from .toml failed!%s", err)
		return nil, err
	}
	//Logger = logging.GetLogger(Conf.Log.Level)

	return Conf, nil
}