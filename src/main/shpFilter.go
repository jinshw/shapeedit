package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/jonas-p/go-shp"
	"shapetest/tools/gis"
)

func main() {
	src := "./shape/Rwulumuqi.shp"
	var dist string = "./shape/Rwulumuqi_gd.shp"
	shape, err := shp.Open(src)
	if err != nil {
		fmt.Errorf("UpdateShape err: %v", err)
	}
	defer shape.Close()

	// fields to write
	fields := shape.Fields()
	shapeNew, errNew := shp.Create(dist, shape.GeometryType)
	if errNew != nil {
		fmt.Errorf("UpdateShape err: %v", errNew)
	}
	defer shapeNew.Close()
	shapeNew.SetFields(fields)

	rowNum := 0
	for shape.Next() {
		n, p := shape.Shape()
		val := shape.ReadAttribute(n, 3)
		kind := val[0:4]
		if kind == "0002" {
			//Box 转化
			box := p.(*shp.PolyLine).Box
			box.MaxX, box.MaxY = gis.GCJ02toWGS84(box.MaxX, box.MaxY)
			box.MinX, box.MinY = gis.GCJ02toWGS84(box.MinX, box.MinY)
			p.(*shp.PolyLine).Box = box

			// Points 转化
			points := p.(*shp.PolyLine).Points
			for i := 0; i < len(points); i++ {
				point := points[i]
				point.X,point.Y = gis.GCJ02toWGS84(point.X,point.Y)
				points[i] = point
			}
			shapeNew.Write(p)
			fmt.Println("n==", n, "    val==", val, "--", val[0:4], "----", p)
			for k, f := range fields {
				val := shape.ReadAttribute(n, k)
				dec := mahonia.NewDecoder("utf8")
				val = dec.ConvertString(val)
				shapeNew.WriteAttribute(rowNum, k, val)
				fmt.Println(k, f)
			}
			rowNum = rowNum + 1
		}

	}

	shape.Close()
	shapeNew.Close()
	fmt.Println("执行完毕！")
}
