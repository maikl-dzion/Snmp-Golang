package main

import (
	"fmt"
	"time"

	// "log"
	"strconv"
)

func load(ch chan string, name string) {
	c := 0
	if name == "" {
		name = "ping"
	}
	for i := 0; ; i++ { c++

		_ , s := DateTypeConvert(c, "")
		line := name + "-" + s

		ch <-line
	}
}
func print(ch chan string) {
	for {
		msg := <-ch
		fmt.Println(msg)
		time.Sleep(time.Millisecond + 1)
	}
}



func main() {

	var forever chan string = make(chan string)


	go load(forever, "ping")
	go print(forever)


	go load(forever, "maikl")
	go print(forever)


	go load(forever, "ini")
	go print(forever)


	var input string
	fmt.Scanln(&input)
}


func DateTypeConvert(num int, str string) (int, string) {

	if num != 0 {
		s := strconv.Itoa(num)
		return num, s
	} else {

		if str == "" {
			return num, str
		}

		d, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println(err)
		}
		return d, str
	}
}