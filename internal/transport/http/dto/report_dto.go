package dto

import (
	"time"

	"github.com/google/uuid"
	"medical-visit/internal/domain/model"
)

type SaveReportRequest struct {
	ConversationSummary  string   `json:"conversation_summary" binding:"required,max=10000"`
	DoctorFeedback       string   `json:"doctor_feedback" binding:"max=10000"`
	MaterialsDistributed bool     `json:"materials_distributed"`
	MaterialIDs          []string `json:"material_ids"`
	AdditionalNotes      string   `json:"additional_notes" binding:"max=10000"`
}

type ReportResponse struct {
	ID                   string    `json:"id"`
	VisitID              string    `json:"visit_id"`
	ConversationSummary  string    `json:"conversation_summary"`
	DoctorFeedback       string    `json:"doctor_feedback,omitempty"`
	MaterialsDistributed bool      `json:"materials_distributed"`
	MaterialIDs          []string  `json:"material_ids,omitempty"`
	AdditionalNotes      string    `json:"additional_notes,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func ParseUUIDs(values []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, len(values))
	for i, value := range values {
		id, err := uuid.Parse(value)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func ToReportResponse(report *model.CallReport) *ReportResponse {
	if report == nil {
		return nil
	}
	materialIDs := make([]string, 0, len(report.Materials))
	for _, material := range report.Materials {
		if material != nil {
			materialIDs = append(materialIDs, material.ID.String())
		}
	}
	return &ReportResponse{ID: report.ID.String(), VisitID: report.VisitID.String(), ConversationSummary: report.DiscussionSummary, DoctorFeedback: report.HCPFeedback, MaterialsDistributed: report.MaterialsDistributed, MaterialIDs: materialIDs, AdditionalNotes: report.AdditionalNotes, CreatedAt: report.CreatedAt, UpdatedAt: report.UpdatedAt}
}
