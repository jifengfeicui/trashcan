package utils

import "math"

const (
	// EarthRadius 地球半径（公里）
	EarthRadius = 6371.0
)

// CalculateDistance 使用Haversine公式计算两点间距离（单位：公里）
// lat1, lng1: 第一个点的纬度和经度
// lat2, lng2: 第二个点的纬度和经度
// 返回：两点间的距离（公里）
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// 将角度转换为弧度
	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	deltaLat := (lat2 - lat1) * math.Pi / 180.0
	deltaLng := (lng2 - lng1) * math.Pi / 180.0

	// Haversine公式
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 计算距离
	distance := EarthRadius * c

	return distance
}
