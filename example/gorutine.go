package main

import (
	"fmt"
	"time"
)

func main() {

	start := time.Now()

	go foo(100000099, "r1")
	go foo(34, "r2")
	go foo(100000099, "r3")
	go foo(100000099, "r4")
	go foo(100000099, "r5")

	elapsedTime := time.Since(start)

	fmt.Println("Total Time For Execution: " + elapsedTime.String())

	time.Sleep(time.Second)

}


func foo(limit int, name string) {

	var sum int = 0
	for i:=0; i < limit; i++ {
		sum += i
	}


	fmt.Println(name)
	fmt.Println(sum)

}

