package main

import "time"

func main() {
	now := time.Now()
	println(now.Format("01-02"))
}
