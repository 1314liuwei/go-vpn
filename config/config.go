package config

import "github.com/spf13/viper"

type Config struct {
	Conn *Conn
	TUN  *TUN
}
type TUN struct {
	Name string
	MTU  int
	Addr string
}
type Conn struct {
	Addr string
	Port int
}

func New() (*Config, error) {
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
