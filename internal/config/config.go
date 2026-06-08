// Package config loads YAML configuration.
package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent struct {
		TickIntervalSeconds   int    `yaml:"tick_interval_seconds"`
		RetentionDays         int    `yaml:"retention_days"`
		IncidentReportDir     string `yaml:"incident_report_dir"`
		MarkdownReportDir     string `yaml:"markdown_report_dir"`
	} `yaml:"agent"`
	Detection struct {
		ErrorSpike struct {
			WindowMinutes          int     `yaml:"window_minutes"`
			BaselineWindowMinutes  int     `yaml:"baseline_window_minutes"`
			StdDevMultiplier       float64 `yaml:"std_dev_multiplier"`
		} `yaml:"error_spike"`
		LatencyDegrade struct {
			WindowMinutes          int     `yaml:"window_minutes"`
			BaselineWindowMinutes  int     `yaml:"baseline_window_minutes"`
			ThresholdMultiplier    float64 `yaml:"threshold_multiplier"`
		} `yaml:"latency_degrade"`
		Heartbeat struct {
			TimeoutSeconds  int            `yaml:"timeout_seconds"`
			ServiceTimeouts map[string]int `yaml:"service_timeouts"`
		} `yaml:"heartbeat"`
	} `yaml:"detection"`
	LLM struct {
		Model              string `yaml:"model"`
		MaxTokensPerCall   int    `yaml:"max_tokens_per_call"`
		CacheTTLHours      int    `yaml:"cache_ttl_hours"`
		TokenLimitTotal    int    `yaml:"token_limit_total"`
		APITimeoutSeconds  int    `yaml:"api_timeout_seconds"`
	} `yaml:"llm"`
	Services []struct {
		Name string `yaml:"name"`
	} `yaml:"services"`
	Storage struct {
		SQLitePath   string `yaml:"sqlite_path"`
		MaxOpenConns int    `yaml:"max_open_conns"`
	} `yaml:"storage"`
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
	// Set defaults if missing
	if cfg.Agent.TickIntervalSeconds == 0 {
		cfg.Agent.TickIntervalSeconds = 60
	}
	if cfg.Agent.RetentionDays == 0 {
		cfg.Agent.RetentionDays = 7
	}
	if cfg.Agent.IncidentReportDir == "" {
		cfg.Agent.IncidentReportDir = "/data/incidents"
	}
	if cfg.Agent.MarkdownReportDir == "" {
		cfg.Agent.MarkdownReportDir = "/data/reports"
	}
	if cfg.LLM.Model == "" {
		cfg.LLM.Model = "deepseek-chat"
	}
	if cfg.LLM.MaxTokensPerCall == 0 {
		cfg.LLM.MaxTokensPerCall = 1000
	}
	if cfg.LLM.CacheTTLHours == 0 {
		cfg.LLM.CacheTTLHours = 24
	}
	if cfg.LLM.TokenLimitTotal == 0 {
		cfg.LLM.TokenLimitTotal = 4500000
	}
	if cfg.LLM.APITimeoutSeconds == 0 {
		cfg.LLM.APITimeoutSeconds = 30
	}
	if cfg.Storage.SQLitePath == "" {
		cfg.Storage.SQLitePath = "/data/relay.db"
	}
	if cfg.Storage.MaxOpenConns == 0 {
		cfg.Storage.MaxOpenConns = 1
	}
	return &cfg, nil
}
