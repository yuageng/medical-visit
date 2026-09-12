package compliance

import (
	"math"
	"testing"
)

func TestCoordinate_Validate(t *testing.T) {
	tests := []struct {
		name      string
		coord     Coordinate
		wantError bool
	}{
		{
			name:      "valid coordinate",
			coord:     Coordinate{Latitude: 31.2304, Longitude: 121.4737},
			wantError: false,
		},
		{
			name:      "latitude too small",
			coord:     Coordinate{Latitude: -91, Longitude: 0},
			wantError: true,
		},
		{
			name:      "latitude too large",
			coord:     Coordinate{Latitude: 91, Longitude: 0},
			wantError: true,
		},
		{
			name:      "longitude too small",
			coord:     Coordinate{Latitude: 0, Longitude: -181},
			wantError: true,
		},
		{
			name:      "longitude too large",
			coord:     Coordinate{Latitude: 0, Longitude: 181},
			wantError: true,
		},
		{
			name:      "NaN latitude",
			coord:     Coordinate{Latitude: math.NaN(), Longitude: 0},
			wantError: true,
		},
		{
			name:      "Inf longitude",
			coord:     Coordinate{Latitude: 0, Longitude: math.Inf(1)},
			wantError: true,
		},
		{
			name:      "boundary latitude -90",
			coord:     Coordinate{Latitude: -90, Longitude: 0},
			wantError: false,
		},
		{
			name:      "boundary latitude 90",
			coord:     Coordinate{Latitude: 90, Longitude: 0},
			wantError: false,
		},
		{
			name:      "boundary longitude -180",
			coord:     Coordinate{Latitude: 0, Longitude: -180},
			wantError: false,
		},
		{
			name:      "boundary longitude 180",
			coord:     Coordinate{Latitude: 0, Longitude: 180},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.coord.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Coordinate.Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestHaversineDistanceMeters(t *testing.T) {
	tests := []struct {
		name          string
		coord1        Coordinate
		coord2        Coordinate
		wantDistance  float64
		tolerance     float64
		wantError     bool
		errorContains string
	}{
		{
			name:         "same coordinate",
			coord1:       Coordinate{Latitude: 31.2304, Longitude: 121.4737},
			coord2:       Coordinate{Latitude: 31.2304, Longitude: 121.4737},
			wantDistance: 0,
			tolerance:    0.1,
			wantError:    false,
		},
		{
			name:         "approximately 500 meters north",
			coord1:       Coordinate{Latitude: 31.2304, Longitude: 121.4737},
			coord2:       Coordinate{Latitude: 31.2349, Longitude: 121.4737},
			wantDistance: 500,
			tolerance:    10,
			wantError:    false,
		},
		{
			name:         "approximately 500 meters east",
			coord1:       Coordinate{Latitude: 31.2304, Longitude: 121.4737},
			coord2:       Coordinate{Latitude: 31.2304, Longitude: 121.4792},
			wantDistance: 523,
			tolerance:    10,
			wantError:    false,
		},
		{
			name:         "cross equator",
			coord1:       Coordinate{Latitude: -1, Longitude: 0},
			coord2:       Coordinate{Latitude: 1, Longitude: 0},
			wantDistance: 222390,
			tolerance:    100,
			wantError:    false,
		},
		{
			name:         "cross prime meridian",
			coord1:       Coordinate{Latitude: 0, Longitude: -1},
			coord2:       Coordinate{Latitude: 0, Longitude: 1},
			wantDistance: 222390,
			tolerance:    100,
			wantError:    false,
		},
		{
			name:          "invalid coord1 latitude",
			coord1:        Coordinate{Latitude: 91, Longitude: 0},
			coord2:        Coordinate{Latitude: 0, Longitude: 0},
			wantError:     true,
			errorContains: "coord1 invalid",
		},
		{
			name:          "invalid coord2 longitude",
			coord1:        Coordinate{Latitude: 0, Longitude: 0},
			coord2:        Coordinate{Latitude: 0, Longitude: 181},
			wantError:     true,
			errorContains: "coord2 invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance, err := HaversineDistanceMeters(tt.coord1, tt.coord2)

			if tt.wantError {
				if err == nil {
					t.Error("HaversineDistanceMeters() expected error but got nil")
					return
				}
				if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("HaversineDistanceMeters() error = %v, want error containing %s", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("HaversineDistanceMeters() unexpected error = %v", err)
				return
			}

			if math.Abs(distance-tt.wantDistance) > tt.tolerance {
				t.Errorf("HaversineDistanceMeters() = %.2f meters, want %.2f ± %.2f meters", distance, tt.wantDistance, tt.tolerance)
			}
		})
	}
}

func TestHaversineDistanceMeters_Boundary(t *testing.T) {
	// 测试恰好 500 米边界
	coord1 := Coordinate{Latitude: 31.2304, Longitude: 121.4737}
	coord2 := Coordinate{Latitude: 31.2349, Longitude: 121.4737}

	distance, err := HaversineDistanceMeters(coord1, coord2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// 验证距离接近 500 米（允许 10 米误差）
	if math.Abs(distance-500) > 10 {
		t.Errorf("Distance = %.2f, expected ~500 meters", distance)
	}

	// 验证比较逻辑
	maxDistance := 500.0
	if distance <= maxDistance {
		t.Log("Distance ≤ 500m: compliant ✓")
	} else {
		t.Log("Distance > 500m: non-compliant")
	}
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
