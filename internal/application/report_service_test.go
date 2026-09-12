package application

import (
	"errors"
	"testing"
	"time"

	"medical-visit/internal/domain/model"
	"medical-visit/internal/infrastructure/clock"

	"github.com/google/uuid"
)

func TestSaveReport_CompletesCheckedOutVisit(t *testing.T) {
	visitID := uuid.New()
	materialID := uuid.New()
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	visit := &model.Visit{ID: visitID, Status: model.VisitStatusCheckedOut}
	var savedReport *model.CallReport
	var savedMaterialIDs []uuid.UUID

	visitRepo := &mockVisitRepository{
		findByIDFunc: func(id uuid.UUID) (*model.Visit, error) {
			return visit, nil
		},
		saveCallReportFunc: func(id uuid.UUID, report *model.CallReport, materialIDs []uuid.UUID) error {
			savedReport = report
			savedMaterialIDs = materialIDs
			visit.Status = model.VisitStatusCompleted
			visit.CallReport = report
			return nil
		},
	}
	masterRepo := &mockMasterDataRepository{
		findAcademicMaterialsByIDsFunc: func(ids []uuid.UUID) ([]*model.AcademicMaterial, error) {
			return []*model.AcademicMaterial{{ID: materialID, Active: true}}, nil
		},
	}
	service := NewReportService(visitRepo, masterRepo, clock.FixedClock{Time: now})

	report, err := service.SaveReport(SaveReportInput{
		VisitID: visitID, ConversationSummary: "介绍临床研究结果",
		DoctorFeedback: "关注长期安全性", MaterialsDistributed: true,
		MaterialIDs: []uuid.UUID{materialID}, AdditionalNotes: "后续跟进",
	})
	if err != nil {
		t.Fatalf("SaveReport() error = %v", err)
	}
	if report == nil || report.ID == uuid.Nil {
		t.Fatal("report was not returned")
	}
	if savedReport == nil || savedReport.DiscussionSummary != "介绍临床研究结果" {
		t.Fatalf("saved report = %+v", savedReport)
	}
	if len(savedMaterialIDs) != 1 || savedMaterialIDs[0] != materialID {
		t.Fatalf("material IDs = %v", savedMaterialIDs)
	}
	if !savedReport.CreatedAt.Equal(now) || !savedReport.UpdatedAt.Equal(now) {
		t.Fatalf("timestamps = %v / %v, want %v", savedReport.CreatedAt, savedReport.UpdatedAt, now)
	}
	if visit.Status != model.VisitStatusCompleted {
		t.Fatalf("visit status = %s, want COMPLETED", visit.Status)
	}
}

func TestSaveReport_RejectsBeforeCheckout(t *testing.T) {
	visitID := uuid.New()
	visitRepo := &mockVisitRepository{
		findByIDFunc: func(id uuid.UUID) (*model.Visit, error) {
			return &model.Visit{ID: visitID, Status: model.VisitStatusCheckedIn}, nil
		},
	}
	service := NewReportService(visitRepo, &mockMasterDataRepository{}, clock.FixedClock{Time: time.Now()})

	_, err := service.SaveReport(SaveReportInput{VisitID: visitID, ConversationSummary: "summary"})
	if !errors.Is(err, ErrReportNotAllowed) {
		t.Fatalf("error = %v, want ErrReportNotAllowed", err)
	}
}

func TestSaveReport_RejectsMissingMaterial(t *testing.T) {
	visitID := uuid.New()
	missingID := uuid.New()
	visitRepo := &mockVisitRepository{
		findByIDFunc: func(id uuid.UUID) (*model.Visit, error) {
			return &model.Visit{ID: visitID, Status: model.VisitStatusCheckedOut}, nil
		},
	}
	masterRepo := &mockMasterDataRepository{
		findAcademicMaterialsByIDsFunc: func(ids []uuid.UUID) ([]*model.AcademicMaterial, error) {
			return []*model.AcademicMaterial{}, nil
		},
	}
	service := NewReportService(visitRepo, masterRepo, clock.FixedClock{Time: time.Now()})

	_, err := service.SaveReport(SaveReportInput{
		VisitID: visitID, ConversationSummary: "summary",
		MaterialsDistributed: true, MaterialIDs: []uuid.UUID{missingID},
	})
	if !errors.Is(err, ErrMaterialNotFound) {
		t.Fatalf("error = %v, want ErrMaterialNotFound", err)
	}
}

func TestSaveReport_AllowsUpdateAfterCompletion(t *testing.T) {
	visitID := uuid.New()
	now := time.Date(2026, 9, 12, 11, 0, 0, 0, time.UTC)
	visit := &model.Visit{ID: visitID, Status: model.VisitStatusCompleted}
	var saveCount int
	visitRepo := &mockVisitRepository{
		findByIDFunc: func(id uuid.UUID) (*model.Visit, error) { return visit, nil },
		saveCallReportFunc: func(id uuid.UUID, report *model.CallReport, materialIDs []uuid.UUID) error {
			saveCount++
			visit.CallReport = report
			return nil
		},
	}
	masterRepo := &mockMasterDataRepository{
		findAcademicMaterialsByIDsFunc: func(ids []uuid.UUID) ([]*model.AcademicMaterial, error) {
			return []*model.AcademicMaterial{}, nil
		},
	}
	service := NewReportService(visitRepo, masterRepo, clock.FixedClock{Time: now})

	_, err := service.SaveReport(SaveReportInput{VisitID: visitID, ConversationSummary: "第一次内容"})
	if err != nil {
		t.Fatalf("first SaveReport() error = %v", err)
	}
	_, err = service.SaveReport(SaveReportInput{VisitID: visitID, ConversationSummary: "修改后的内容"})
	if err != nil {
		t.Fatalf("second SaveReport() error = %v", err)
	}
	if saveCount != 2 {
		t.Fatalf("save count = %d, want 2", saveCount)
	}
	if visit.CallReport.DiscussionSummary != "修改后的内容" {
		t.Fatalf("summary = %q, want updated value", visit.CallReport.DiscussionSummary)
	}
}

func TestSaveReport_RejectsMaterialIDsWhenNotDistributed(t *testing.T) {
	visitID := uuid.New()
	materialID := uuid.New()
	visitRepo := &mockVisitRepository{
		findByIDFunc: func(id uuid.UUID) (*model.Visit, error) {
			return &model.Visit{ID: visitID, Status: model.VisitStatusCheckedOut}, nil
		},
	}
	service := NewReportService(visitRepo, &mockMasterDataRepository{}, clock.FixedClock{Time: time.Now()})

	_, err := service.SaveReport(SaveReportInput{
		VisitID: visitID, ConversationSummary: "summary",
		MaterialIDs: []uuid.UUID{materialID},
	})
	if err == nil {
		t.Fatal("SaveReport() expected error, got nil")
	}
}
