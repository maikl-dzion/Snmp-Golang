package main

import "fmt"



func changeValue(x *int){
	*x = 56 + 7
}


func main() {

	d := 5
	fmt.Println("d before:", d)		// 5
	changeValue(&d)					// изменяем значение
	fmt.Println("d after:", d)		// 25 - значение изменилось!
}