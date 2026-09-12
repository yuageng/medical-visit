package application

import (
	"errors"
	"fmt"
	"strings"

	"medical-visit/internal/domain/model"
	"medical-visit/internal/domain/repository"
	"medical-visit/internal/infrastructure/clock"

	"github.com/google/uuid"
)

var (
	ErrReportNotAllowed            = errors.New("visit status does not allow call report")
	ErrConversationSummaryRequired = errors.New("conversation summary is required")
	ErrMaterialNotFound            = errors.New("academic material not found")
	ErrMaterialNotActive           = errors.New("academic material is not active")
)

type SaveReportInput struct {
	VisitID              uuid.UUID
	ConversationSummary  string
	DoctorFeedback       string
	MaterialsDistributed bool
	MaterialIDs          []uuid.UUID
	AdditionalNotes      string
}

type ReportService struct {
	visitRepo  repository.VisitRepository
	masterRepo repository.MasterDataRepository
	clock      clock.Clock
}

func NewReportService(visitRepo repository.VisitRepository, masterRepo repository.MasterDataRepository, clk clock.Clock) *ReportService {
	return &ReportService{visitRepo: visitRepo, masterRepo: masterRepo, clock: clk}
}

func (s *ReportService) SaveReport(input SaveReportInput) (*model.CallReport, error) {
	visit, err := s.visitRepo.FindByID(input.VisitID)
	if err != nil {
		return nil, fmt.Errorf("find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}
	if visit.Status != model.VisitStatusCheckedOut && visit.Status != model.VisitStatusCompleted {
		return nil, ErrReportNotAllowed
	}
	if strings.TrimSpace(input.ConversationSummary) == "" {
		return nil, ErrConversationSummaryRequired
	}
	if !input.MaterialsDistributed && len(input.MaterialIDs) > 0 {
		return nil, errors.New("material_ids require materials_distributed=true")
	}
	materials, err := s.masterRepo.FindAcademicMaterialsByIDs(input.MaterialIDs)
	if err != nil {
		return nil, fmt.Errorf("find academic materials: %w", err)
	}
	if len(materials) != len(input.MaterialIDs) {
		return nil, ErrMaterialNotFound
	}
	for _, material := range materials {
		if !material.Active {
			return nil, ErrMaterialNotActive
		}
	}
	now := s.clock.Now().UTC()
	report := &model.CallReport{
		ID: uuid.New(), VisitID: input.VisitID, DiscussionSummary: strings.TrimSpace(input.ConversationSummary),
		HCPFeedback: input.DoctorFeedback, MaterialsDistributed: input.MaterialsDistributed,
		AdditionalNotes: input.AdditionalNotes, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.visitRepo.SaveCallReport(input.VisitID, report, input.MaterialIDs); err != nil {
		return nil, fmt.Errorf("save call report: %w", err)
	}
	visit, err = s.visitRepo.FindByID(input.VisitID)
	if err != nil {
		return nil, fmt.Errorf("reload visit: %w", err)
	}
	if visit != nil && visit.CallReport != nil {
		return visit.CallReport, nil
	}
	return report, nil
}
