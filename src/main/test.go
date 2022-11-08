package main

import "fmt"

func main() {

	score := 100.1238456
	f_score := float64(score)
	s_score := fmt.Sprintf("%.3f", f_score)
	fmt.Println(s_score)
}
