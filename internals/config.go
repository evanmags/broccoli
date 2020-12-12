package internals

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	PluginsDir string `yaml:"plugins_dir"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
