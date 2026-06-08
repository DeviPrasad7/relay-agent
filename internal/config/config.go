// Package config loads YAML configuration.
//
// Last updated: 2026-06-09
package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Detection struct {
		ErrorSpike struct {
			WindowMinutes      int     `yaml:"window_minutes"`
			BaselineWindowMinutes int  `yaml:"baseline_window_minutes"`
			StdDevMultiplier   float64 `yaml:"std_dev_multiplier"`
		} `yaml:"error_spike"`
	} `yaml:"detection"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
