/**
 * @Author: Administrator
 * @Description:
 * @File:  utils
 * @Version: 1.0.0
 * @Date: 2019/12/16 16:02
 */

package main

import (
	"fmt"
	"github.com/jonas-p/go-shp"
)

func GetCols(shapefile string) []shp.Field {
	shape, err := shp.Open(shapefile)
	if err != nil {
		fmt.Errorf("UpdateShape err: %v", err)
	}
	defer shape.Close()
	fields := shape.Fields()
	return fields
}
