package config

import "github.com/spf13/viper"

type (
    //Config - interface for getting at configuration
	Config interface {
		GetOpsManagerHost() string
		GetOpsManagerUser() string
		GetOpsManagerPassword() string
		GetAdminUser() string
		GetAdminPassword() string
		GetString(key string) (value string, err error)
		GetBool(key string) (value bool, err error)
		GetInt(key string) (value int, err error)
		GetSubConfig(key string) (config Config, err error)
		GetKeys() (keys []string)
	}
    //FileConfig - struct for File based config
	FileConfig struct {
		Viper viper.Viper
	}
)
