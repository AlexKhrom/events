package main

import (
	"fmt"
	"strconv"
)

func main() {
	str := `"some str"`
	fmt.Println("hello world!!! = ", strconv.Quote(str)[1:len(strconv.Quote(str))-1])
	fmt.Println("hello world!!! = ", str)

}

//package main
//
//import (
//	"fmt"
//	"math"
//)
//
//func main() {
//	var f1x1, f1x2, f1y1, f1y2, f2x1, f2x2, f2y1, f2y2 int
//	fmt.Scanf("%d %d %d %d", &f1x1, &f1y1, &f1x2, &f1y2)
//	fmt.Scanf("%d %d %d %d", &f2x1, &f2y1, &f2x2, &f2y2)
//
//	var lenF, widthF float64
//	if f1x1 < f2x2 {
//		lenF = math.Abs(float64(f1x1) - float64(f2x2))
//	} else {
//		lenF = math.Abs(float64(f2x1) - float64(f1x2))
//	}
//	if f1y1 < f2y2 {
//		widthF = math.Abs(float64(f1y1) - float64(f2y2))
//	} else {
//		widthF = math.Abs(float64(f2y1) - float64(f1y2))
//	}
//
//	if lenF > widthF {
//		fmt.Println(lenF * lenF)
//	} else {
//		fmt.Println(widthF * widthF)
//	}
//
//}
