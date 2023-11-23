package main

import (
	"fmt"
	"time"
)

var PATH = "data/test_data_1/"

func main() {
	var n int
	fmt.Print("Введите количество горутин: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		fmt.Println("Ошибка при считывании числа:", err)
		return
	}

	kc := NewKeywordCounter()

	startTime := time.Now()

	resultMap, err := kc.Count(n, PATH+"input", PATH+"keywords")
	if err != nil {
		fmt.Println("Ошибка во времени выполнения: ", err)
	}
	endTime := time.Now()

	fmt.Println(resultMap)
	fmt.Println("Время исполнения: ", endTime.Sub(startTime))
}
