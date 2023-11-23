package main

import "fmt"

func main() {

	result, err := Unpack("a4bc2d5e2")

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print(result)
}
