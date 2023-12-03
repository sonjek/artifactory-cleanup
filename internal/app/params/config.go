package params

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	Type            string   `yaml:"type"`
	Repos           []string `yaml:"repos"`
	CleanupPatterns []string `yaml:"cleanupPatterns"`
	ExcludePatterns []string `yaml:"excludePatterns"`
	Rules           []Rule   `yaml:"rule"`
}

type Rule struct {
	Name    string `yaml:"name"`
	Count   int    `yaml:"count"`
	Pattern string `yaml:"pattern"`
}

func ParceConfigFile(configFile string) (*ConfigFile, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return &ConfigFile{}, err
	}

	var config ConfigFile

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return &ConfigFile{}, err
	}

	if len(config.Rules) > 1 {
		return &ConfigFile{}, fmt.Errorf("found more then one rule")
	}

	return &config, nil
}
