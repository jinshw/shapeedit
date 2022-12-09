package main

import (
	"bufio"
	list2 "container/list"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	path := "tj.txt"
	//path := "D:\\项目\\广州先导区项目\\车辆轨迹.txt"
	// 1、读取原始车辆轨迹数据
	txt := ReadTxt(path)
	fmt.Println(txt)
	var jsonList []interface{}
	err := json.Unmarshal([]byte(txt), &jsonList)
	if err != nil {
		fmt.Println(err)
	}

	// 2、在原始两点间插值，达到每秒30
	var startLon float64 = 0.0
	var startLat float64 = 0.0
	var endLon float64 = 0.0
	var endLat float64 = 0.0
	var startMap map[string]interface{}
	//var endMap map[string]interface{}

	var stepNum float64 = 60

	resultList := list2.New()

	for index, lineJson := range jsonList {
		fmt.Println(index, lineJson)
		lineMap := lineJson.(map[string]interface{})
		if index == 0 {
			//startLon = lineMap["longitude"].(float64)
			//startLat = lineMap["latitude"].(float64)
			//startMap = lineMap
			//resultList.PushBack(lineMap)
		} else {
			startLon = startMap["longitude"].(float64)
			startLat = startMap["latitude"].(float64)
			endLon = lineMap["longitude"].(float64)
			endLat = lineMap["latitude"].(float64)

			distance := GetDistanceTwo(startLon, startLat, endLon, endLat)
			azimuth := ComputeAzimuth(startLon, startLat, endLon, endLat)
			//distanceInt := distance
			stepLen := distance / stepNum

			for i := 0; i < int(stepNum); i++ {
				lineInsertMap := make(map[string]interface{})
				stepLenKM := stepLen
				//stepLenKM :=  stepLen * float64(i+1)
				lonlat := ConvertDistanceToLogLat(startLon, startLat, stepLenKM, azimuth)
				lonlatArr := strings.Split(lonlat, ",")
				lon, _ := strconv.ParseFloat(lonlatArr[0], 64)
				lat, _ := strconv.ParseFloat(lonlatArr[1], 64)
				startLon = lon
				startLat = lat
				vid := lineMap["vid"]
				vlicense := lineMap["vlicense"]
				updated_at := lineMap["updated_at"]
				speed := lineMap["speed"]
				lineInsertMap["vid"] = vid
				lineInsertMap["vlicense"] = vlicense
				lineInsertMap["updated_at"] = updated_at
				lineInsertMap["longitude"] = lon
				lineInsertMap["latitude"] = lat
				lineInsertMap["heading"] = azimuth
				lineInsertMap["speed"] = speed
				resultList.PushBack(lineInsertMap)
				lineMap["heading"] = azimuth

			}
		}
		resultList.PushBack(lineMap)
		startMap = lineMap
	}

	fmt.Println(resultList)

	fileCSV, err := os.OpenFile("./out_car_line_tj.csv", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// 关闭文件
	defer fileCSV.Close()
	count := 0
	lineStr := ""
	for i := resultList.Front(); i != nil; i = i.Next() {
		lineObj := i.Value.(map[string]interface{})
		fmt.Println("front--->back:", count, i.Value)
		lon := strconv.FormatFloat(lineObj["longitude"].(float64), 'f', -1, 64)
		lat := strconv.FormatFloat(lineObj["latitude"].(float64), 'f', -1, 64)
		heading := strconv.FormatFloat(lineObj["heading"].(float64), 'f', -1, 64)
		speed := strconv.FormatFloat(lineObj["speed"].(float64), 'f', -1, 64)
		vid := lineObj["vid"].(string)
		vlicense := lineObj["vlicense"].(string)
		updated_at := lineObj["updated_at"].(string)
		//lineStr = lineObj["vid"].(string) + "," + lineObj["vlicense"].(string) + "," +
		//	lineObj["updated_at"].(string) + "," +
		//	heading + "," +
		//	lon + "," +
		//	lat + "," +
		//	speed + ","
		//lineStr = "[" + lon + "," + lat + "]"+","
		lineStr = vid + "," + vlicense + "," + speed + "," + updated_at + "," + lon + "," + lat + "," + heading
		fileCSV.WriteString(lineStr + "\r\n")
		count++
	}

}

func ReadTxt(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	resultStr := ""
	for scanner.Scan() {
		resultStr = resultStr + scanner.Text()
	}
	return resultStr
}
