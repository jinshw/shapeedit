package main

import (
	list2 "container/list"
	"fmt"
	"github.com/jinshw/go-shp"
	"log"
	"os"
	"strconv"
	"strings"
)

/**
* 获取路网桩号和经纬度对应数据
 */
func main() {

	//lon1 := 92.4819020463424
	//lat1 := 43.2509190028098
	//lon2 := 92.4820159670633
	//lat2 := 43.2509174343921
	//distance := GetDistanceTwo(lon1, lat1, lon2, lat2)
	//azimuth := GetAzimuth(lon1, lat1, lon2, lat2)
	//fmt.Println("distance==", distance)
	//fmt.Println("azimuth==", azimuth)

	var stepLen int32 = 1
	stepFloat := float64(stepLen)
	stepLenKM := stepFloat / 1000.0

	reader, err := shp.Open("C:\\Users\\DELL\\Desktop\\shape\\GL_GSGL_LD_441.shp")
	list := list2.New()

	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	//fields := reader.Fields()
	var startPoint shp.Point
	var endPoint shp.Point
	for reader.Next() {
		_, p := reader.Shape()
		polyLine := p.(*shp.PolyLine)
		points := polyLine.Points
		fmt.Println(len(points))
		roadcode := reader.GetValue("ROADCODE")
		ldlx := reader.GetValue("LDLX")
		ldlx = strings.Split(ldlx, ".")[0]
		startzh_m, _ := strconv.ParseFloat(reader.GetValue("STARTZH_M"), 64)
		end_m, _ := strconv.ParseFloat(reader.GetValue("ENDZH_M"), 64)
		zhFlag := end_m - startzh_m
		metaId := reader.GetValue("META_ID")

		var distanceSum float64 = 0.0

		for i := 1; i < len(points); i++ {

			//point := points[i]
			//fmt.Println(i, roadcode, point.X, point.Y, startzh_m, end_m, metaId, n)
			startPoint = points[i-1]
			endPoint = points[i]

			distance := GetDistanceTwo(startPoint.X, startPoint.Y, endPoint.X, endPoint.Y)
			azimuth := ComputeAzimuth(startPoint.X, startPoint.Y, endPoint.X, endPoint.Y)
			//azimuth := GetAzimuth(startPoint.X, startPoint.Y, endPoint.X, endPoint.Y)
			distanceInt := int32(distance * 1000)
			if distanceInt > stepLen {
				step := int(distanceInt / stepLen)
				for idx := 0; idx < step; idx++ {
					lonlat := ConvertDistanceToLogLat(startPoint.X, startPoint.Y, stepLenKM, azimuth)
					lonlatArr := strings.Split(lonlat, ",")
					startPoint.X, _ = strconv.ParseFloat(lonlatArr[0], 64)
					startPoint.Y, _ = strconv.ParseFloat(lonlatArr[1], 64)
					var zh float64
					if zhFlag > 0 {
						//zhSum = zhSum + float64(int32((idx+1))*stepLen)
						zh = startzh_m + float64(int32((idx+1))*stepLen) + (distanceSum * 1000)
					} else {
						//zhSum = zhSum - float64(int32((idx+1))*stepLen)
						zh = startzh_m - float64(int32((idx+1))*stepLen) - (distanceSum * 1000)
					}
					//zh = startzh_m + zhSum
					line := roadcode + "," + lonlat + "," + fmt.Sprintf("%.2f", zh) + "," + metaId + "," + ldlx
					//line := roadcode + "," + lonlat + "," + strconv.FormatFloat(zh, 'f', 30, 64) + ","
					list.PushBack(line)
				}
			} else {
				lonlat := ConvertDistanceToLogLat(startPoint.X, startPoint.Y, stepLenKM, azimuth)
				lonlatArr := strings.Split(lonlat, ",")
				startPoint.X, _ = strconv.ParseFloat(lonlatArr[0], 64)
				startPoint.Y, _ = strconv.ParseFloat(lonlatArr[1], 64)
				var zh float64
				if zhFlag > 0 {
					//zhSum = zhSum + float64(stepLen)
					zh = startzh_m + float64(stepLen) + (distanceSum * 1000)
				} else {
					//zhSum = zhSum - float64(stepLen)
					zh = startzh_m - float64(stepLen) - (distanceSum * 1000)
				}
				//zh = startzh_m + zhSum
				line := roadcode + "," + lonlat + "," + fmt.Sprintf("%.2f", zh) + "," + metaId + "," + ldlx
				//line := roadcode + "," + lonlat + "," + strconv.FormatFloat(zh, 'f', 30, 64) + "," + metaId
				list.PushBack(line)
			}
			distanceSum = distanceSum + distance
			//if i == 100 {
			//	break
			//}
		}
		// print feature
		//fmt.Println(reflect.TypeOf(p).Elem(), p.BBox())
		// print attributes
		//for k, f := range fields {
		//	val := reader.ReadAttribute(n, k)
		//	fmt.Printf("\t%v: %v\n", f, val)
		//}
		//fmt.Println()
	}

	file, err := os.OpenFile("D:/helloworld.csv", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// 关闭文件
	defer file.Close()
	count := 0
	for i := list.Front(); i != nil; i = i.Next() {
		fmt.Println("front--->back:", count, i.Value)
		file.WriteString(i.Value.(string) + "\r\n")
		count++
	}

}
