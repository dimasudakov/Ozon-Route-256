package main

import (
	"fmt"
	"math"
)

type Apple struct {
	Id int
}

type Box struct {
	Id     int
	Apples []Apple
}

type Car struct {
	Id    int
	Boxes []Box
}

func putApples(apples []Apple, boxCapacity int) []Car {
	cars := make([]Car, 2)
	for i := 0; i < 2; i++ {
		cars[i].Id = i + 1
	}

	boxesCnt := (len(apples) + boxCapacity - 1) / boxCapacity

	for i := 1; i <= boxesCnt; i++ {
		var box Box
		box.Id = i
		for len(box.Apples) < boxCapacity && len(apples) > 0 {
			box.Apples = append(box.Apples, apples[len(apples)-1])
			apples = apples[:len(apples)-1]
		}
		if i%2 == 0 {
			cars[1].Boxes = append(cars[1].Boxes, box)
		} else {
			cars[0].Boxes = append(cars[0].Boxes, box)
		}
	}

	return cars
}

func showCars(cars []Car) {
	for _, car := range cars {
		for _, box := range car.Boxes {
			fmt.Printf("Машина: %d, Ящик: %d, Яблоки: %v\n", car.Id, box.Id, box.Apples)
		}
	}
}

func main() {

	var applesCnt = int(math.Pow(10, 2))

	apples := make([]Apple, applesCnt)
	for i := 0; i < applesCnt; i++ {
		apples[i].Id = i + 1
	}

	cars := putApples(apples, 10)
	showCars(cars)
}
