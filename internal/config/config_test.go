package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("SERVER_PORT", "8081")
	os.Setenv("COMPLIANCE_MIN_DURATION_SECONDS", "600")
	os.Setenv("COMPLIANCE_MAX_DISTANCE_METERS", "1000")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("COMPLIANCE_MIN_DURATION_SECONDS")
		os.Unsetenv("COMPLIANCE_MAX_DISTANCE_METERS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Database.Host != "testhost" {
		t.Errorf("Database.Host = %v, want testhost", cfg.Database.Host)
	}

	if cfg.Database.Port != 5433 {
		t.Errorf("Database.Port = %v, want 5433", cfg.Database.Port)
	}

	if cfg.Server.Port != 8081 {
		t.Errorf("Server.Port = %v, want 8081", cfg.Server.Port)
	}

	if cfg.Compliance.MinDurationSeconds != 600 {
		t.Errorf("Compliance.MinDurationSeconds = %v, want 600", cfg.Compliance.MinDurationSeconds)
	}

	if cfg.Compliance.MaxDistanceMeters != 1000 {
		t.Errorf("Compliance.MaxDistanceMeters = %v, want 1000", cfg.Compliance.MaxDistanceMeters)
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		DBName:   "dbname",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=user password=pass dbname=dbname sslmode=disable"
	if dsn := cfg.DSN(); dsn != expected {
		t.Errorf("DSN() = %v, want %v", dsn, expected)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{
					Host: "0.0.0.0",
					Port: 8080,
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "pass",
					DBName:   "dbname",
					SSLMode:  "disable",
				},
				Compliance: ComplianceConfig{
					MinDurationSeconds: 300,
					MaxDistanceMeters:  500,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid server port",
			config: Config{
				Server: ServerConfig{
					Host: "0.0.0.0",
					Port: 0,
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "pass",
					DBName:   "dbname",
				},
				Compliance: ComplianceConfig{
					MinDurationSeconds: 300,
					MaxDistanceMeters:  500,
				},
			},
			wantErr: true,
		},
		{
			name: "negative compliance values",
			config: Config{
				Server: ServerConfig{
					Host: "0.0.0.0",
					Port: 8080,
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "pass",
					DBName:   "dbname",
				},
				Compliance: ComplianceConfig{
					MinDurationSeconds: -1,
					MaxDistanceMeters:  500,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
