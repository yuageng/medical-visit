package handler

import (
	"errors"
	"log"

	"medical-visit/internal/application"
	"medical-visit/internal/transport/http/dto"
	"medical-visit/internal/transport/http/response"

	"github.com/gin-gonic/gin"
)

// VisitHandler 拜访 HTTP 处理器
type VisitHandler struct {
	visitService  *application.VisitService
	reportService *application.ReportService
}

// NewVisitHandler 创建拜访处理器实例
func NewVisitHandler(visitService *application.VisitService, reportServices ...*application.ReportService) *VisitHandler {
	var reportService *application.ReportService
	if len(reportServices) > 0 {
		reportService = reportServices[0]
	}
	return &VisitHandler{visitService: visitService, reportService: reportService}
}

// CreateVisit 创建拜访
// POST /api/v1/visits
func (h *VisitHandler) CreateVisit(c *gin.Context) {
	var req dto.CreateVisitRequest

	// 1. 绑定和校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数不合法", err.Error())
		return
	}

	// 2. 解析 UUID
	mrID, err := dto.ParseUUID(req.MRID)
	if err != nil {
		response.BadRequest(c, "MR ID 格式不正确", err.Error())
		return
	}

	hcpID, err := dto.ParseUUID(req.HCPID)
	if err != nil {
		response.BadRequest(c, "HCP ID 格式不正确", err.Error())
		return
	}

	hospitalID, err := dto.ParseUUID(req.HospitalID)
	if err != nil {
		response.BadRequest(c, "Hospital ID 格式不正确", err.Error())
		return
	}

	departmentID, err := dto.ParseUUID(req.DepartmentID)
	if err != nil {
		response.BadRequest(c, "Department ID 格式不正确", err.Error())
		return
	}

	productID, err := dto.ParseUUID(req.ProductID)
	if err != nil {
		response.BadRequest(c, "Product ID 格式不正确", err.Error())
		return
	}

	// 3. 调用应用服务
	input := application.CreateVisitInput{
		MRID:           mrID,
		HCPID:          hcpID,
		HospitalID:     hospitalID,
		DepartmentID:   departmentID,
		ProductID:      productID,
		PlannedStartAt: req.PlannedStartAt,
		PlanNote:       req.PlanNote,
	}

	visit, err := h.visitService.CreateVisit(input)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// 4. 返回响应
	response.Created(c, dto.ToVisitResponse(visit))
}

// CheckIn 签到拜访。
// POST /api/v1/visits/:id/check-in
func (h *VisitHandler) CheckIn(c *gin.Context) {
	visitID, err := dto.ParseUUID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Visit ID 格式不正确", err.Error())
		return
	}

	var req dto.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "签到参数不合法", err.Error())
		return
	}

	visit, err := h.visitService.CheckIn(application.CheckInInput{
		VisitID:     visitID,
		Latitude:    *req.Latitude,
		Longitude:   *req.Longitude,
		CheckInTime: req.CheckInTime,
	})
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	response.Success(c, dto.ToVisitResponse(visit))
}

// CheckOut 签退拜访。
// POST /api/v1/visits/:id/check-out
func (h *VisitHandler) CheckOut(c *gin.Context) {
	visitID, err := dto.ParseUUID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Visit ID 格式不正确", err.Error())
		return
	}
	var req dto.CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "签退参数不合法", err.Error())
		return
	}
	visit, err := h.visitService.CheckOut(application.CheckOutInput{
		VisitID: visitID, Latitude: *req.Latitude, Longitude: *req.Longitude, CheckOutTime: req.CheckOutTime,
	})
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	response.Success(c, dto.ToVisitResponse(visit))
}

// GetVisit 获取拜访详情
// GET /api/v1/visits/:id
func (h *VisitHandler) GetVisit(c *gin.Context) {
	// 1. 解析路径参数
	idStr := c.Param("id")
	id, err := dto.ParseUUID(idStr)
	if err != nil {
		response.BadRequest(c, "Visit ID 格式不正确", err.Error())
		return
	}

	// 2. 调用应用服务
	visit, err := h.visitService.GetVisitByID(id)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// 3. 返回响应
	response.Success(c, dto.ToVisitResponse(visit))
}

// SaveReport 创建或修改拜访记录。
// POST /api/visits/:id/report
func (h *VisitHandler) SaveReport(c *gin.Context) {
	if h.reportService == nil {
		response.InternalServerError(c, "拜访记录服务未初始化")
		return
	}
	visitID, err := dto.ParseUUID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Visit ID 格式不正确", err.Error())
		return
	}
	var req dto.SaveReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "拜访记录参数不合法", err.Error())
		return
	}
	materialIDs, err := dto.ParseUUIDs(req.MaterialIDs)
	if err != nil {
		response.BadRequest(c, "material_id 格式不正确", err.Error())
		return
	}
	report, err := h.reportService.SaveReport(application.SaveReportInput{
		VisitID:              visitID,
		ConversationSummary:  req.ConversationSummary,
		DoctorFeedback:       req.DoctorFeedback,
		MaterialsDistributed: req.MaterialsDistributed,
		MaterialIDs:          materialIDs,
		AdditionalNotes:      req.AdditionalNotes,
	})
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	response.Success(c, dto.ToReportResponse(report))
}

// handleServiceError 处理服务层错误并映射为 HTTP 响应
func (h *VisitHandler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, application.ErrReportNotAllowed):
		response.Conflict(c, "REPORT_NOT_ALLOWED", "当前拜访状态不允许填写拜访记录")
	case errors.Is(err, application.ErrReportConflict):
		response.Conflict(c, "REPORT_CONFLICT", "拜访记录已被其他请求修改，请刷新后重试")
	case errors.Is(err, application.ErrConversationSummaryRequired):
		response.BadRequest(c, "谈话要点不能为空", nil)
	case errors.Is(err, application.ErrMaterialNotFound):
		response.UnprocessableEntity(c, "存在不存在或无效的学术资料", nil)
	case errors.Is(err, application.ErrMaterialNotActive):
		response.UnprocessableEntity(c, "学术资料已停用", nil)
	case errors.Is(err, application.ErrDuplicateMaterialID):
		response.BadRequest(c, "material_id 不允许重复", nil)
	case errors.Is(err, application.ErrMRNotFound):
		response.NotFound(c, "医药代表不存在")
	case errors.Is(err, application.ErrHCPNotFound):
		response.NotFound(c, "医生不存在")
	case errors.Is(err, application.ErrHospitalNotFound):
		response.NotFound(c, "医院不存在")
	case errors.Is(err, application.ErrDepartmentNotFound):
		response.NotFound(c, "科室不存在")
	case errors.Is(err, application.ErrProductNotFound):
		response.NotFound(c, "产品不存在")
	case errors.Is(err, application.ErrVisitNotFound):
		response.NotFound(c, "拜访不存在")
	case errors.Is(err, application.ErrVisitCannotCheckIn):
		response.Conflict(c, "VISIT_CANNOT_CHECK_IN", "拜访当前状态不允许签到")
	case errors.Is(err, application.ErrVisitCannotCheckOut):
		response.Conflict(c, "VISIT_CANNOT_CHECK_OUT", "拜访当前状态不允许签退")
	case errors.Is(err, application.ErrInvalidCheckOutTime):
		response.UnprocessableEntity(c, "签退时间不能早于签到时间", nil)
	case errors.Is(err, application.ErrInvalidCheckOutGPS):
		response.BadRequest(c, "签退 GPS 坐标不合法", nil)
	case errors.Is(err, application.ErrHospitalCoordinatesMissing):
		response.UnprocessableEntity(c, "目标医院没有完整 GPS 坐标", nil)
	case errors.Is(err, application.ErrInvalidCheckInTime):
		response.BadRequest(c, "签到时间不合法", nil)
	case errors.Is(err, application.ErrInvalidCheckInGPS):
		response.BadRequest(c, "签到 GPS 坐标不合法", nil)
	case errors.Is(err, application.ErrDepartmentNotBelongToHospital):
		response.UnprocessableEntity(c, "科室不属于该医院", nil)
	case errors.Is(err, application.ErrHCPNotInHospital):
		response.UnprocessableEntity(c, "医生不在该医院", nil)
	case errors.Is(err, application.ErrHCPNotInDepartment):
		response.UnprocessableEntity(c, "医生不在该科室", nil)
	case errors.Is(err, application.ErrPlannedTimeInPast):
		response.UnprocessableEntity(c, "计划时间不能早于当前时间", nil)
	default:
		log.Printf("Unexpected service error: %v", err)
		response.InternalServerError(c, "服务器内部错误")
	}
}
