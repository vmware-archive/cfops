package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	configName         = "config"
	opsManagerHost     = "opsmanagerhost"
	opsManagerUser     = "opsmanageruser"
	opsManagerPassword = "opsmanagerpass"
	adminUser          = "adminuser"
	adminPassword      = "adminpass"
)

//GetOpsManagerHost - returns the operations manager host name
func (s *FileConfig) GetOpsManagerHost() string {
	return s.Viper.GetString(opsManagerHost)
}

//GetOpsManagerUser - returns the operations manager user name
func (s *FileConfig) GetOpsManagerUser() string {
	return s.Viper.GetString(opsManagerUser)
}

//GetOpsManagerPassword - returns the operations manager password
func (s *FileConfig) GetOpsManagerPassword() string {
	return s.Viper.GetString(opsManagerPassword)
}

//GetAdminUser - returns the admin user name
func (s *FileConfig) GetAdminUser() string {
	return s.Viper.GetString(adminUser)
}

//GetAdminPassword - returns the admin password
func (s *FileConfig) GetAdminPassword() string {
	return s.Viper.GetString(adminPassword)
}

//GetString - returns string value from config based on key, error is key doesn't exist
func (s *FileConfig) GetString(key string) (value string, err error) {
	if s.Viper.InConfig(strings.ToLower(key)) {
		value = s.Viper.GetString(key)
	} else {
		err = fmt.Errorf("key %s not found as a string, valid keys are %s", key, s.Viper.AllKeys())
	}
	return
}

//GetBool - returns boolean value from config based on key, error is key doesn't exist
func (s *FileConfig) GetBool(key string) (value bool, err error) {
	if s.Viper.InConfig(strings.ToLower(key)) {
		value = s.Viper.GetBool(key)
	} else {
		err = fmt.Errorf("key %s not found as a string, valid keys are %s", key, s.Viper.AllKeys())
	}
	return
}

//GetInt - returns int value from config based on key, error is key doesn't exist
func (s *FileConfig) GetInt(key string) (value int, err error) {
	if s.Viper.InConfig(strings.ToLower(key)) {
		value = s.Viper.GetInt(key)
	} else {
		err = fmt.Errorf("key %s not found as a string, valid keys are %s", key, s.Viper.AllKeys())
	}
	return
}

//GetKeys - returns []string of keys
func (s *FileConfig) GetKeys() (keys []string) {
	keys = s.Viper.AllKeys()
	return
}

//GetSubConfig - returns a Config object for a nested level
func (s *FileConfig) GetSubConfig(key string) (config Config, err error) {
	v := s.Viper.Sub(key)
	config = &FileConfig{
		Viper: *v,
	}
	return
}

var NewConfig = func(path string) (config Config, err error) {

	v := viper.New()
	v.SetConfigName(configName)
	v.AddConfigPath(path)

	err = v.ReadInConfig()

	if err != nil {
		return
	}

	config = &FileConfig{
		Viper: *v,
	}
	return
}
