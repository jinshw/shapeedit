package main

import (
	"fmt"
	"github.com/jonas-p/go-shp"
)

func main() {
	addFields := []shp.Field{
		shp.StringField("name12", 64),
		shp.StringField("name13", 64),
	}
	shp.UpdateShape("test.shp", addFields, func(shape *shp.Reader, shapeNew *shp.Writer) {
		fieldsObj := shape.Fields()
		for shape.Next() {
			n, p := shape.Shape()
			fmt.Println(n)
			shapeNew.Write(p)

			// print attributes
			count := 0
			for k, f := range fieldsObj {
				val := shape.ReadAttribute(n, k)
				name := shape.GetValue("NAME")
				fmt.Println("NAME === " + name)
				shapeNew.WriteAttribute(n, k, val)
				count = k
				fmt.Println(f, val)
			}
			fmt.Println(count)
			shapeNew.WriteAttribute(n, count+1, "燕山国际")
			shapeNew.WriteAttribute(n, count+2, "燕山国际2")
		}
	})
}
