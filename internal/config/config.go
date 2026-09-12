package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Compliance ComplianceConfig
	Log        LogConfig
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ComplianceConfig 合规规则配置
type ComplianceConfig struct {
	MinDurationSeconds int
	MaxDistanceMeters  float64
}

// LogConfig 日志配置
type LogConfig struct {
	Level string
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "medical_visit"),
			Password: getEnv("DB_PASSWORD", "medical_visit_password"),
			DBName:   getEnv("DB_NAME", "medical_visit"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Compliance: ComplianceConfig{
			MinDurationSeconds: getEnvAsInt("COMPLIANCE_MIN_DURATION_SECONDS", 300),
			MaxDistanceMeters:  getEnvAsFloat("COMPLIANCE_MAX_DISTANCE_METERS", 500.0),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate 校验配置合法性
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}

	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	if c.Compliance.MinDurationSeconds < 0 {
		return fmt.Errorf("min duration seconds must be non-negative")
	}

	if c.Compliance.MaxDistanceMeters < 0 {
		return fmt.Errorf("max distance meters must be non-negative")
	}

	return nil
}

// DSN 返回数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取整型环境变量
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsFloat 获取浮点型环境变量
func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}
