package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type ServeConfig struct {
	Port int `yaml:"port"`
}

type ManageConfig struct {
	Port               int           `yaml:"port"`
	Buckets            []string      `yaml:"buckets"`
	DefaultKeyLifetime time.Duration `yaml:"defaultKeyLifetime"`
}

type AppConfig struct {
	FilePath           string       `yaml:"filePath"`
	Host               string       `yaml:"host"`
	DBFile             string       `yaml:"dbFile"`
	HostKeys           []string     `yaml:"hostKeys"`
	AuthorizedKeyFiles []string     `yaml:"authorizedKeys"`
	Serve              ServeConfig  `yaml:"serve"`
	Manage             ManageConfig `yaml:"manage"`
}

type Metadata struct {
	Version string `yaml:"version"`
}

type Config struct {
	App  AppConfig `yaml:"app"`
	Meta Metadata  `yaml:"meta"`
}

func Load(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
