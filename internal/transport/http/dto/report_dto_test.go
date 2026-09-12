package dto

import (
	"testing"

	"medical-visit/internal/domain/model"

	"github.com/google/uuid"
)

func TestToReportResponse_IncludesMaterialIDs(t *testing.T) {
	firstID, secondID := uuid.New(), uuid.New()
	response := ToReportResponse(&model.CallReport{
		ID: uuid.New(), VisitID: uuid.New(), DiscussionSummary: "summary",
		Materials: []*model.AcademicMaterial{{ID: firstID}, nil, {ID: secondID}},
	})
	if response == nil || len(response.MaterialIDs) != 2 {
		t.Fatalf("material IDs = %v", response.MaterialIDs)
	}
	if response.MaterialIDs[0] != firstID.String() || response.MaterialIDs[1] != secondID.String() {
		t.Fatalf("material IDs = %v", response.MaterialIDs)
	}
}
