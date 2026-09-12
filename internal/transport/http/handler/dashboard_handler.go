package handler

import (
	"net/http"

	"medical-visit/internal/application"
	"medical-visit/internal/transport/http/dto"
	"medical-visit/internal/transport/http/response"

	"github.com/gin-gonic/gin"
)

// DashboardHandler Dashboard HTTP 处理器。
type DashboardHandler struct {
	service *application.DashboardService
}

func NewDashboardHandler(service *application.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// MonthlyVisitsByProduct GET /api/dashboard/visits/monthly?month=2026-09
func (h *DashboardHandler) MonthlyVisitsByProduct(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		response.BadRequest(c, "month 参数不能为空", nil)
		return
	}
	items, err := h.service.MonthlyVisitsByProduct(month)
	if err != nil {
		response.BadRequest(c, "month 参数格式应为 YYYY-MM", err.Error())
		return
	}
	c.JSON(http.StatusOK, response.Response{
		Code: "SUCCESS", Message: "操作成功",
		Data: dto.MonthlyVisitsResponse{Month: month, Items: items},
	})
}
