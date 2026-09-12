package main

import (
	"fmt"
	"log"
	"os"

	"medical-visit/internal/config"
	"medical-visit/internal/domain/compliance"
	"medical-visit/internal/infrastructure/persistence"
	"medical-visit/internal/transport/http"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting medical-visit API server on %s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Database: %s@%s:%d/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Printf("Compliance rules: min_duration=%ds, max_distance=%.1fm", cfg.Compliance.MinDurationSeconds, cfg.Compliance.MaxDistanceMeters)

	// 连接数据库
	db, err := persistence.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 开发环境下自动迁移表结构（生产环境应使用 migrate 工具）
	if os.Getenv("AUTO_MIGRATE") == "true" {
		log.Println("Running auto migration...")
		if err := persistence.AutoMigrate(db); err != nil {
			log.Fatalf("Failed to auto migrate: %v", err)
		}
		log.Println("Auto migration completed")
	}

	// 初始化路由（带依赖注入），传入环境配置的合规阈值。
	rules := compliance.Rules{MinDurationSeconds: cfg.Compliance.MinDurationSeconds, MaxDistanceMeters: cfg.Compliance.MaxDistanceMeters}
	router := http.SetupRouterWithDependencies(db, rules)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server is listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
