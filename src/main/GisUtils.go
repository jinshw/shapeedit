package main

import (
	"math"
	"strconv"
)

var EARTH_RADIUS float64 = 6378.137
var EARTH_ARC float64 = 111.199

func GetDistanceTwo(lon1 float64, lat1 float64, lon2 float64, lat2 float64) float64 {
	radLat1 := rad(lat1)
	radLat2 := rad(lat2)
	a := radLat1 - radLat2
	b := rad(lon1) - rad(lon2)
	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * EARTH_RADIUS
	return s
}

func GetAzimuth(lon1 float64, lat1 float64, lon2 float64, lat2 float64) float64 {
	lat1 = rad(lat1)
	lat2 = rad(lat2)
	lon1 = rad(lon1)
	lon2 = rad(lon2)
	var azimuth float64
	azimuth = math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1)
	azimuth = math.Sqrt(1 - azimuth*azimuth)
	azimuth = math.Cos(lat2) * math.Sin(lon2-lon1) / azimuth
	azimuth = math.Asin(azimuth) * 180 / math.Pi

	//if lon1 < lon2 {
	//	azimuth = 90.0
	//} else {
	//	azimuth = 270.0
	//}
	return azimuth

}

/**
/*
     * 获取两个经纬度坐标点的角度
     * @param LatLng
     * @param LatLng
*/
// 计算方位角,正北向为0度，以顺时针方向递增
func ComputeAzimuth(lon1 float64, lat1 float64, lon2 float64, lat2 float64) float64 {
	var ilat1 int64 = int64(0.50 + lat1*360000.0)
	var ilat2 int64 = (int64)(0.50 + lat2*360000.0)
	var ilon1 int64 = (int64)(0.50 + lon1*360000.0)
	var ilon2 int64 = (int64)(0.50 + lon2*360000.0)

	lat1 = rad(lat1)
	lon1 = rad(lon1)
	lat2 = rad(lat2)
	lon2 = rad(lon2)

	var result float64 = 0.0

	if (ilat1 == ilat2) && (ilon1 == ilon2) {
		return result
	} else if ilon1 == ilon2 {
		if ilat1 > ilat2 {
			result = 180.0
		}

	} else {
		var c float64 = math.Acos(math.Sin(lat2)*math.Sin(lat1) + math.Cos(lat2)*math.Cos(lat1)*math.Cos((lon2-lon1)))
		aSinVal := math.Cos(lat2) * math.Sin((lon2 - lon1)) / math.Sin(c)
		if aSinVal > 1 {
			aSinVal = 1
		}
		if aSinVal < -1 {
			aSinVal = -1
		}
		var A float64 = math.Asin(aSinVal)
		result = toDegress(A)
		if (ilat2 > ilat1) && (ilon2 > ilon1) {
			return result
		} else if (ilat2 < ilat1) && (ilon2 < ilon1) {
			result = 180.0 - result
		} else if (ilat2 < ilat1) && (ilon2 > ilon1) {
			result = 180.0 - result
		} else if (ilat2 > ilat1) && (ilon2 < ilon1) {
			result += 360.0
		}
	}
	return result

}

func ConvertDistanceToLogLat(lng1 float64, lat1 float64, distance float64, azimuth float64) string {
	azimuth = rad(azimuth)
	// 将距离转换成经度的计算公式
	lon := lng1 + (distance*math.Sin(azimuth))/(EARTH_ARC*math.Cos(rad(lat1)))
	// 将距离转换成纬度的计算公式
	lat := lat1 + (distance*math.Cos(azimuth))/EARTH_ARC
	return strconv.FormatFloat(lon, 'f', 30, 64) + "," + strconv.FormatFloat(lat, 'f', 30, 64)
}

func rad(d float64) float64 {
	return d * math.Pi / 180.0
}
func toDegress(r float64) float64 {
	return r * 180.0 / math.Pi
}
