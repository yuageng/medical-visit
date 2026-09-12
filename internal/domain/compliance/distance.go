package compliance

import (
	"errors"
	"fmt"
	"math"
)

// Coordinate 地理坐标值对象
type Coordinate struct {
	Latitude  float64 // 纬度，范围 [-90, 90]
	Longitude float64 // 经度，范围 [-180, 180]
}

// NewCoordinate 创建坐标并校验
func NewCoordinate(latitude, longitude float64) (Coordinate, error) {
	coord := Coordinate{
		Latitude:  latitude,
		Longitude: longitude,
	}
	if err := coord.Validate(); err != nil {
		return Coordinate{}, err
	}
	return coord, nil
}

// Validate 校验坐标合法性
func (c Coordinate) Validate() error {
	if math.IsNaN(c.Latitude) || math.IsInf(c.Latitude, 0) {
		return errors.New("latitude must be a finite number")
	}
	if c.Latitude < -90 || c.Latitude > 90 {
		return fmt.Errorf("latitude must be in [-90, 90], got %f", c.Latitude)
	}
	if math.IsNaN(c.Longitude) || math.IsInf(c.Longitude, 0) {
		return errors.New("longitude must be a finite number")
	}
	if c.Longitude < -180 || c.Longitude > 180 {
		return fmt.Errorf("longitude must be in [-180, 180], got %f", c.Longitude)
	}
	return nil
}

// IsZero 判断是否为零值
func (c Coordinate) IsZero() bool {
	return c.Latitude == 0 && c.Longitude == 0
}

// HaversineDistanceMeters 计算两点之间的 Haversine 距离（米）
// 使用球面大圆距离公式，适用于地球表面短距离计算
func HaversineDistanceMeters(coord1, coord2 Coordinate) (float64, error) {
	// 校验坐标合法性
	if err := coord1.Validate(); err != nil {
		return 0, fmt.Errorf("coord1 invalid: %w", err)
	}
	if err := coord2.Validate(); err != nil {
		return 0, fmt.Errorf("coord2 invalid: %w", err)
	}

	// 地球平均半径（米）
	const earthRadiusMeters = 6371000.0

	// 转换为弧度
	lat1Rad := coord1.Latitude * math.Pi / 180
	lat2Rad := coord2.Latitude * math.Pi / 180
	deltaLatRad := (coord2.Latitude - coord1.Latitude) * math.Pi / 180
	deltaLonRad := (coord2.Longitude - coord1.Longitude) * math.Pi / 180

	// Haversine 公式
	a := math.Sin(deltaLatRad/2)*math.Sin(deltaLatRad/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLonRad/2)*math.Sin(deltaLonRad/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadiusMeters * c
	return distance, nil
}
