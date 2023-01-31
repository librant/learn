package main

import (
	"fmt"
	"math/rand"
)

func main() {
	numbers := []byte("12345")
	letters := []byte("ABCDE")

	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
		letters[i], letters[j] = letters[j], letters[i]
	})
	for i := range numbers {
		fmt.Printf("%c: %c\n", letters[i], numbers[i])
	}
}
