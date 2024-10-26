package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/yazmeyaa/duck_wars/pkg/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database Database `yaml:"database"`
	JWT      JWT      `yaml:"jwt"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  bool   `yaml:"sslmode"`
}

type JWT struct {
	Secret string `yaml:"secret"`
}

func (d *Database) DSN() (string, error) {
	if d.Username == "" || d.Database == "" {
		return "", fmt.Errorf("username and database are required fields")
	}

	var sb strings.Builder

	sb.WriteString("postgresql://")
	sb.WriteString(d.Username)
	sb.WriteString(":")
	sb.WriteString(d.Password)
	sb.WriteString("@")

	if d.Host == "" {
		sb.WriteString("localhost")
	} else {
		sb.WriteString(d.Host)
	}

	if d.Port != 0 {
		sb.WriteString(fmt.Sprintf(":%d", d.Port))
	} else {
		sb.WriteString(":5432")
	}

	sb.WriteString("/")
	sb.WriteString(d.Database)

	sb.WriteString(fmt.Sprintf("?sslmode=%s", boolToSSLMode(d.SSLMode)))

	return sb.String(), nil
}

func boolToSSLMode(sslMode bool) string {
	if sslMode {
		return "enable"
	}
	return "disable"
}

func New(log *slog.Logger) (*Config, error) {
	const funcName string = "config.New"
	config := Config{}

	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		logger.LogError(log, funcName, "Error while reading config file", err)
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		logger.LogError(log, funcName, "Error while unmarshal config file", err)
		return nil, err
	}

	return &config, nil
}
