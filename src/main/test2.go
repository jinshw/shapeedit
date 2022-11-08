package main

import (
	"fmt"
	"github.com/jinshw/go-shp"
)

func main() {
	var startPoint shp.Point
	var endPoint shp.Point
	var startPoint2 shp.Point
	var endPoint2 shp.Point
	var stepLenKM = 0.01
	startPoint.X = 87.5987666305756
	startPoint.Y = 44.0893549975352
	endPoint.X = 87.6009050529259
	endPoint.Y = 44.0893549975352

	startPoint2.X = 87.6025214664925
	startPoint2.Y = 44.0893549975352
	endPoint2.X = 87.6045039832723
	endPoint2.Y = 44.0891284241375

	//distance := GetDistanceTwo(startPoint.X, startPoint.Y, endPoint.X, endPoint.Y)
	azimuth := GetAzimuth(startPoint.X, startPoint.Y, endPoint.X, endPoint.Y)
	azimuth2 := GetAzimuth(startPoint2.X, startPoint2.Y, endPoint2.X, endPoint2.Y)

	azimuth3 := ComputeAzimuth(startPoint2.X, startPoint2.Y, endPoint2.X, endPoint2.Y)
	azimuth2 = 180 - azimuth2
	lonlat := ConvertDistanceToLogLat(startPoint.X, startPoint.Y, stepLenKM, azimuth)
	lonlat2 := ConvertDistanceToLogLat(startPoint2.X, startPoint2.Y, stepLenKM, azimuth2)

	fmt.Println("azimuth=", azimuth)
	fmt.Println("azimuth2=", azimuth2)
	fmt.Println("azimuth3=", azimuth3)
	fmt.Println("lonlat=", lonlat)
	fmt.Println("lonlat2=", lonlat2)

}
