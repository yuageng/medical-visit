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
	ErrReportConflict              = errors.New("call report was concurrently modified")
	ErrConversationSummaryRequired = errors.New("conversation summary is required")
	ErrMaterialNotFound            = errors.New("academic material not found")
	ErrMaterialNotActive           = errors.New("academic material is not active")
	ErrDuplicateMaterialID         = errors.New("duplicate academic material id")
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
	seenMaterialIDs := make(map[uuid.UUID]struct{}, len(input.MaterialIDs))
	for _, materialID := range input.MaterialIDs {
		if _, exists := seenMaterialIDs[materialID]; exists {
			return nil, ErrDuplicateMaterialID
		}
		seenMaterialIDs[materialID] = struct{}{}
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
	if err := s.visitRepo.SaveCallReport(input.VisitID, visit.Version, report, input.MaterialIDs); err != nil {
		if errors.Is(err, repository.ErrVisitStateConflict) {
			return nil, ErrReportConflict
		}
		if errors.Is(err, repository.ErrMaterialStateConflict) {
			return nil, ErrMaterialNotActive
		}
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
