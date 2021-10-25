package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	DB     DBConf
	// TODO
}

type DBConf struct {
	url      string
	inMemory bool
}

type LoggerConf struct {
	File  string
	Level string
	// TODO
}

func NewConfig(configPath string) Config {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config file: %w", err))
	}
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("decode into struct: %w", err))
	}

	return c
}

// TODO
