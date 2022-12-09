package main

import (
	"bufio"
	list2 "container/list"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	path := "D:\\项目\\广州先导区项目\\车辆轨迹.txt"
	// 1、读取原始车辆轨迹数据
	txt := readTxt(path)
	fmt.Println(txt)
	var jsonList []interface{}
	err := json.Unmarshal([]byte(txt), &jsonList)
	if err != nil {
		fmt.Println(err)
	}

	// 2、在原始两点间插值，达到每秒30


	resultList := list2.New()


	for index, lineJson := range jsonList {
		fmt.Println(index, lineJson)
		lineMap := lineJson.(map[string]interface{})
		lineInsertMap := make(map[string]interface{})
		vid := lineMap["vid"]
		vlicense := lineMap["vlicense"]
		updated_at := lineMap["updated_at"]
		speed := lineMap["speed"]
		lon := lineMap["longitude"]
		lat := lineMap["latitude"]
		heading := lineMap["heading"]
		lineInsertMap["vid"] = vid
		lineInsertMap["vlicense"] = vlicense
		lineInsertMap["updated_at"] = updated_at
		lineInsertMap["longitude"] = lon
		lineInsertMap["latitude"] = lat
		lineInsertMap["heading"] = heading
		lineInsertMap["speed"] = speed
		resultList.PushBack(lineInsertMap)
	}


	fileCSV, err := os.OpenFile("./out_car_old.csv", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
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
		//heading := strconv.FormatFloat(lineObj["heading"].(float64), 'f', -1, 64)
		//speed := strconv.FormatFloat(lineObj["speed"].(float64), 'f', -1, 64)
		//lineStr = lineObj["vid"].(string) + "," + lineObj["vlicense"].(string) + "," +
		//	lineObj["updated_at"].(string) + "," +
		//	heading + "," +
		//	lon + "," +
		//	lat + "," +
		//	speed + ","
		lineStr = "[" + lon + "," + lat + "]"+","
		fileCSV.WriteString(lineStr + "\r\n")
		count++
	}

}

func readTxt(path string) string {
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
