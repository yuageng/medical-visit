package http

import (
	"medical-visit/internal/application"
	"medical-visit/internal/domain/compliance"
	"medical-visit/internal/infrastructure/clock"
	"medical-visit/internal/infrastructure/persistence"
	"medical-visit/internal/transport/http/handler"
	"medical-visit/internal/transport/http/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 初始化不依赖数据库的基础路由。
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := newRouter()
	router.GET("/healthz", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "healthy"})
	})
	router.GET("/api/v1/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"message": "pong"})
	})
	return router
}

// SetupRouterWithDependencies 初始化完整 API 路由。
func SetupRouterWithDependencies(db *gorm.DB, complianceRules ...compliance.Rules) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := newRouter()
	router.GET("/healthz", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "healthy"})
	})
	router.GET("/api/v1/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"message": "pong"})
	})

	visitRepo := persistence.NewVisitRepository(db)
	masterDataRepo := persistence.NewMasterDataRepository(db)
	clk := clock.RealClock{}
	rules := compliance.DefaultRules()
	if len(complianceRules) > 0 {
		rules = complianceRules[0]
	}
	complianceService, err := compliance.NewService(rules)
	if err != nil {
		panic(err)
	}
	visitService := application.NewVisitService(visitRepo, masterDataRepo, clk, complianceService)
	reportService := application.NewReportService(visitRepo, masterDataRepo, clk)
	dashboardService := application.NewDashboardService(visitRepo)
	visitHandler := handler.NewVisitHandler(visitService, reportService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	legacyVisits := router.Group("/api/visits")
	legacyVisits.POST("/:id/check-in", visitHandler.CheckIn)
	legacyVisits.POST("/:id/check-out", visitHandler.CheckOut)
	legacyVisits.POST("/:id/report", visitHandler.SaveReport)
	legacyDashboard := router.Group("/api/dashboard")
	legacyDashboard.GET("/visits/monthly", dashboardHandler.MonthlyVisitsByProduct)

	v1 := router.Group("/api/v1")
	visits := v1.Group("/visits")
	visits.POST("", visitHandler.CreateVisit)
	visits.GET("/:id", visitHandler.GetVisit)
	visits.POST("/:id/check-in", visitHandler.CheckIn)
	visits.POST("/:id/check-out", visitHandler.CheckOut)
	visits.POST("/:id/report", visitHandler.SaveReport)
	v1Dashboard := v1.Group("/dashboard")
	v1Dashboard.GET("/visits/monthly", dashboardHandler.MonthlyVisitsByProduct)

	return router
}

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(Recovery(), Logger(), RequestID(), CORS())
	return router
}
