package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 调用待测试的函数
	config, err := InitConfigFromConfigFile("test")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// server port should be 8080
	if config.Server.Port != "8080" {
		t.Fatalf("Server port should be 8080")
	}

	// if database is enabled, then database host should be localhost
	if config.Database.Enabled && config.Database.Host != "localhost" {
		t.Fatalf("Database host should be localhost")
	}

	// if database is enabled, port should be 5432
	if config.Database.Enabled && config.Database.Port != "5432" {
		t.Fatalf("Database port should be 5432")
	}
}
