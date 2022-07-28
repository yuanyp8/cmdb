package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

var global *Config

func (c *Config) InitGlobal() error {
	// 后期还会有别的初始化配置，可能会报错，所以这里预留error
	global = c
	return nil
}

// LoadConfigFromToml 从配置文件加载配置
func LoadConfigFromToml(filepath string) (*Config, error) {
	cfg := newConfig()
	_, err := toml.DecodeFile(filepath, cfg)
	if err != nil {
		return nil, fmt.Errorf("load config from file error, path:%s, %s", filepath, err)
	}
	return cfg, nil
}

func LoadConfigFromTomlAndENV(filepath string) (*Config, error) {
	cfg, err := LoadConfigFromToml(filepath)
	if err != nil {
		return nil, err
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("load config from ENV error, %s", err)
	}
	return cfg, nil
}

// C 全局配置对象
func C() *Config {
	if global == nil {
		panic("Load Config first")
	}
	return global
}

func LoadConfig(filepath string) error {
	cfg, err := LoadConfigFromTomlAndENV(filepath)
	if err != nil {
		return err
	}
	if err := cfg.InitGlobal(); err != nil {
		return err
	}
	return nil
}
