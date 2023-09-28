package main

import (
	"fastcsv/splitter"
	"fmt"
	"time"
)

func main() {
	s, err := splitter.Open("./test.csv", "./dist", "best_datetime")
	if err != nil {
		fmt.Println(err)
		return
	}

	before := time.Now()
	if err := s.Split(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Time cost %v\n", time.Since(before).Milliseconds())
}
