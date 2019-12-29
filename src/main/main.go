package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"github.com/jonas-p/go-shp"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Data struct {
	Filename string              `json:"filename"`
	JoinObj  map[string]string   ` json:"joinObj" `
	AddList  []map[string]string `json:"addList" `
}

var Encoding = "utf8"
var ShapeFileName string = "test.shp"
var DirPathPrefix = "./shape/"

var DictData map[string]map[string]string

func FileDownload(c *gin.Context, filename string) {
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filename)
}

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.Static("/index", "./html")

	r.GET("/getCols", func(context *gin.Context) {
		filename := context.Query("filename")
		Encoding = context.Query("encoding")
		filepath := DirPathPrefix + filename
		ShapeFileName = filepath
		type Column struct {
			Id   int    `json:"id"`
			Name string `json："name"`
		}
		type Result struct {
			Code int      `json:"code"`
			Msg  string   `json:"msg"`
			Data []Column `json:"data"`
		}
		results := make([]Column, 0)
		//for i := 0; i < 10; i++ {
		//	var col Column
		//	col.Id = i
		//	col.Name = "列名称" + strconv.Itoa(i)
		//	results = append(results, col)
		//}

		fields := GetCols(filepath)
		for i := range (fields) {
			fmt.Println(i, string(fields[i].Name[:]))
			var col Column
			col.Id = i
			name := string(fields[i].Name[:])
			col.Name = strings.Trim(strings.Replace(name, "\u0000", "", -1), " ")
			results = append(results, col)
		}
		var result Result
		result.Code = 0
		result.Msg = "success"
		result.Data = results
		context.JSON(http.StatusOK, result)
	})

	r.GET("/list", func(c *gin.Context) {

		// open a shapefile for reading
		shape, err := shp.Open("test.shp")
		if err != nil {
			log.Fatal(err)
		}
		defer shape.Close()

		// fields from the attribute table (DBF)
		fields := shape.Fields()

		// loop through all features in the shapefile
		for shape.Next() {
			n, p := shape.Shape()

			// print feature
			fmt.Println(reflect.TypeOf(p).Elem(), p.BBox())

			// print attributes
			for k, f := range fields {
				val := shape.ReadAttribute(n, k)
				fmt.Printf("\t%v: %v\n", f, val)
			}
			if n == 41519 {
				fmt.Println("n===41519")
			}
			fmt.Println(n)
		}

	})

	r.GET("/shpList", func(c *gin.Context) {
		limit := c.Query("limit")
		page := c.Query("page")
		filename := c.Query("filename")
		fmt.Println(filename)
		filepath := DirPathPrefix + filename

		pageInt, _ := strconv.Atoi(page)
		limitInt, _ := strconv.Atoi(limit)
		var startPageNumber int
		var endPageNumber int
		startPageNumber = (pageInt-1)*limitInt + 1
		endPageNumber = pageInt * limitInt

		type Result struct {
			Code  int                 `json:"code"`
			Msg   string              `json:"msg"`
			Count int                 `json:"count"`
			Data  []map[string]string `json:"data"`
		}

		shape, err := shp.Open(filepath)
		if err != nil {
			fmt.Errorf("UpdateShape err: %v", err)
		}
		defer shape.Close()
		fields := shape.Fields()

		results := make([]map[string]string, 0)
		for shape.Next() {
			n, p := shape.Shape()
			fmt.Println(n, p, fields, startPageNumber)
			if n >= (startPageNumber-1) && n <= (endPageNumber-1) {
				var item map[string]string
				item = make(map[string]string)

				for k, f := range fields {
					val := shape.ReadAttribute(n, k)
					//name := shape.GetValue("NAME")
					fmt.Println(n, p, k, f, val)

					name := string(f.Name[:])
					name = strings.Trim(strings.Replace(name, "\u0000", "", -1), " ")
					dec := mahonia.NewDecoder(Encoding)
					valutf8 := dec.ConvertString(val)
					item[name] = valutf8
				}
				results = append(results, item)
			}
			if n >= endPageNumber {
				break
			}

		}

		var result Result
		result.Code = 0
		result.Msg = "success"
		result.Count = shape.AttributeCount()
		result.Data = results

		c.JSON(http.StatusOK, result)
	})

	r.GET("/exports", func(context *gin.Context) {
		shape, err := shp.Open(ShapeFileName)
		if err != nil {
			fmt.Errorf("UpdateShape err: %v", err)
		}
		defer shape.Close()
		fields := shape.Fields()

		file := xlsx.NewFile()
		sheet, _ := file.AddSheet("Sheet1")
		// add Header
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		for k, f := range fields {
			cell := row.AddCell()
			cell.Value = f.String()
			fmt.Println(k)
		}

		for shape.Next() {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度
			n, p := shape.Shape()

			for k, f := range fields {
				val := shape.ReadAttribute(n, k)
				dec := mahonia.NewDecoder(Encoding)
				val = dec.ConvertString(val)

				fmt.Println(n, p, k, f, val)
				cell := row.AddCell()
				// 去除特殊字符
				getVal := strings.Replace(val, " ", "", -1)
				getVal = strings.Replace(strconv.Quote(getVal), "\x00", "", -1)
				getVal = strings.Replace(getVal, "\"", "", -1)
				getVal = strings.Replace(getVal, "\\x00", "", -1)
				cell.Value = getVal
			}

		}

		errXlsx := file.Save("file.xlsx")
		if errXlsx != nil {
			panic(errXlsx)
		}

		FileDownload(context, "file.xlsx")
	})

	r.GET("/getExcel", func(context *gin.Context) {
		pathname := "excel"
		type Result struct {
			Code int                 `json:"code"`
			Msg  string              `json:"msg"`
			Data []map[string]string `json:"data"`
		}
		rd, err := ioutil.ReadDir(pathname)
		if err != nil {
			fmt.Println("读取excel文件夹失败")
		}
		var data = make([]map[string]string, 0)
		for _, fi := range rd {
			item := make(map[string]string)
			if fi.IsDir() {
				fullDir := pathname + "/" + fi.Name()
				fmt.Println("文件夹路径=" + fullDir)
			} else {
				filename := pathname + "/" + fi.Name()
				item["file"] = filename

				fmt.Println("文件路径==" + filename)
			}
			data = append(data, item)
		}
		var result Result
		result.Code = 0
		result.Msg = "success"
		result.Data = data
		context.JSON(http.StatusOK, result)
	})

	r.GET("/getExcelHead", func(context *gin.Context) {
		filename := context.Query("filename")
		fmt.Println("filename==" + filename)
		type Result struct {
			Code int               `json:"code"`
			Msg  string            `json:"msg"`
			Data map[string]string `json:"data"`
		}
		var data = make(map[string]string)
		data = readExcelHead(filename)
		var result Result
		result.Code = 0
		result.Msg = "success"
		result.Data = data
		context.JSON(http.StatusOK, result)
	})

	r.POST("/saveShp", func(context *gin.Context) {
		//type Data struct {
		//	Filename string              `json:"filename"`
		//	JoinObj  map[string]string   ` json:"joinObj" `
		//	AddList  []map[string]string `json:"addList" `
		//}
		var reqData Data

		if err := context.ShouldBindJSON(&reqData); err == nil {
			fmt.Printf("login info:%#v\n", reqData)
		} else {
			log.Println("context.ShouldBindJSON error:" + err.Error())
		}
		result, dict := readExcel(reqData.Filename, reqData)
		DictData = dict
		updateShape(ShapeFileName, reqData, result, dict)
		fmt.Println(result, dict)

		type Result struct {
			Code int               `json:"code"`
			Msg  string            `json:"msg"`
			Data map[string]string `json:"data"`
		}
		var robj Result
		robj.Code = 200
		context.JSON(http.StatusOK, robj)
	})

	r.POST("/getFileList", func(context *gin.Context) {
		//dirPaht := context.PostForm("username")
		type FileList struct {
			FilePath string `json:"filePath"`
		}
		type Result struct {
			Code int      `json:"code"`
			Msg  string   `json:"msg"`
			Data []string `json:"data"`
		}
		files, err := ListDir(DirPathPrefix, ".shp")
		if err != nil {
			fmt.Println(err)
		}

		var result Result
		result.Code = 200
		result.Msg = "success"
		result.Data = files
		context.JSON(http.StatusOK, result)
	})

	r.Run()
}

func readExcelHead(excelPath string) map[string]string {
	var result map[string]string
	result = make(map[string]string)
	excelFileName := excelPath
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Print(err)
	}
	//for _, sheet := range xlFile.Sheets {
	//	for _, row := range sheet.Rows {
	//		for _, cell := range row.Cells {
	//			text := cell.String()
	//			fmt.Printf("%s\n", text)
	//			result[text] = text
	//		}
	//	}
	//}

	sheet := xlFile.Sheets[0]
	row := sheet.Rows[0]
	for _, cell := range row.Cells {
		text := cell.String()
		result[text] = text
	}
	return result
}

func readExcel(excelPath string, obj Data) ([]map[string]string, map[string]map[string]string) {

	var dictKeyStr string
	var dictKeyVal string
	var dictValStr string
	var dictValFiledMap map[string]string
	dictValFiledMap = make(map[string]string)
	//var dictValVal string
	dictKeyStr = "#"
	for key, val := range obj.JoinObj {
		fmt.Println(key, val)
		dictKeyStr = dictKeyStr + val + "#"
	}
	dictValStr = "#"
	for bindKey, bindVal := range obj.AddList {
		fmt.Println(bindKey, bindVal, bindVal["bindFiledText"])
		dictValStr = dictValStr + bindVal["bindFiledText"] + "#"
	}

	var result []map[string]string
	var teampMap map[string]string
	var dict map[string]map[string]string
	dict = make(map[string]map[string]string)
	excelFileName := excelPath
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Print(err)
	}

	rowCount := 0

	for _, sheet := range xlFile.Sheets {

		for _, row := range sheet.Rows {
			dictKeyVal = "#"
			count := 0
			teampMap = make(map[string]string)
			rowMap := make(map[string]string)

			for _, cell := range row.Cells {
				text := cell.String()
				fmt.Printf("%s\n", text)
				teampMap[ strconv.Itoa(count)] = text
				if rowCount == 0 {
					if strings.Contains(dictKeyStr, "#"+text+"#") {
						dictKeyStr = strings.Replace(dictKeyStr, text, strconv.Itoa(count), -1)
					}
					if strings.Contains(dictValStr, "#"+text+"#") {
						dictValStr = strings.Replace(dictValStr, text, strconv.Itoa(count), -1)
						dictValFiledMap[strconv.Itoa(count)] = text
					}

				} else { //生成excel数据字段
					indexs := strings.Split(dictKeyStr, "#")
					for keyIdx := range indexs {
						fmt.Println(keyIdx)
						if indexs[keyIdx] == strconv.Itoa(count) {
							dictKeyVal = dictKeyVal + text + "#"
						}
					}

					if strings.Contains(dictValStr, "#"+strconv.Itoa(count)+"#") {
						rowMap[dictValFiledMap[strconv.Itoa(count)]] = text
					}

				}
				count = count + 1
			}
			result = append(result, teampMap)
			dict[dictKeyVal] = rowMap
			rowCount = rowCount + 1
		}
	}
	return result, dict
}

func updateShape(shapeFile string, obj Data, result []map[string]string, dict map[string]map[string]string) {
	var dictKeyStr string
	addFields := []shp.Field{}
	dictValStr := "#"
	for bindKey, bindVal := range obj.AddList {
		fmt.Println(bindKey, bindVal, bindVal["newFiledText"], bindVal["bindFiledText"])
		newFiledText := bindVal["newFiledText"]
		fmt.Println("newFiledText==" + newFiledText)
		addFields = append(addFields, shp.StringField(newFiledText, 64))
		dictValStr = dictValStr + bindVal["bindFiledText"] + "#"
	}

	shp.UpdateShape(shapeFile, addFields, func(shape *shp.Reader, shapeNew *shp.Writer) {
		fieldsObj := shape.Fields()
		for shape.Next() {
			n, p := shape.Shape()
			shapeNew.Write(p)

			// print attributes
			//count := 0
			dictKeyStr = "#"
			for key, val := range obj.JoinObj {
				fmt.Println(key, val)
				shpVal := shape.GetValue(key)
				getVal := strings.Replace(shpVal, " ", "", -1)
				getVal = strings.Replace(strconv.Quote(getVal), "\x00", "", -1)
				getVal = strings.Replace(getVal, "\"", "", -1)
				getVal = strings.Replace(getVal, "\\x00", "", -1)
				dictKeyStr = dictKeyStr + getVal + "#"
			}

			for k, f := range fieldsObj {
				val := shape.ReadAttribute(n, k)
				dec := mahonia.NewDecoder(Encoding)
				val = dec.ConvertString(val)
				shapeNew.WriteAttribute(n, k, val)
				//keyFiled := f.String()
				fmt.Println(k, f)

				//if strings.Contains(dictValStr, "#"+keyFiled+"#") {
				//	shapeNew.WriteAttribute(n, count, DictData[dictKeyStr][keyFiled])
				//	fmt.Println(dictKeyStr, keyFiled)
				//	fmt.Println(n, count, dict[dictKeyStr][keyFiled])
				//} else {
				//	shapeNew.WriteAttribute(n, k, val)
				//}

				//name := shape.GetValue("NAME")
				//fmt.Println("NAME === " + name)
				//shapeNew.WriteAttribute(n, k, val)
				//count = k
				//fmt.Println(f, shape.Attribute(2))
			}
			teampBindField := dictValStr[1 : len(dictValStr)-1]
			bindFileds := strings.Split(teampBindField, "#")
			for idx := 0; idx < len(bindFileds); idx++ {
				bindFiled := bindFileds[idx]
				newVal := DictData[dictKeyStr][bindFiled]
				//dec := mahonia.NewDecoder(Encoding)
				//newVal = dec.ConvertString(newVal)
				shapeNew.WriteAttribute(n, len(fieldsObj)+idx, newVal)
			}
			//shapeNew.WriteAttribute(n, count+1, 101)
			//shapeNew.WriteAttribute(n, count+2, "男10and女1")
		}
	})
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤
func ListDir(dirPath string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() {
			continue
			//filesTeamp, _ := ListDir(dirPath+PthSep+fi.Name(),suffix)
			//files = append(files,filesTeamp...)
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPath+PthSep+fi.Name())
		}
	}
	return files, nil
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤
func WalkDir(dirPath string, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	err = filepath.Walk(dirPath, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if suffix != "*" {
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
				files = append(files, filename)
			}
		} else {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}
