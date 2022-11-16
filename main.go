package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Filter unique random numbers and print
func printer(numbersGen chan int, count int, quit chan struct{}) {
	numbers := make(map[int]bool)
	defer close(quit)
	defer close(numbersGen)

	for number := range numbersGen {
		if numbers[number] {
			continue
		}
		numbers[number] = true
		fmt.Println(number)
		count--

		if count < 1 {
			return
		}
	}
}

// Worker for generation random numbers to channel
func randomGenWorker(numbersGen chan int, quit chan struct{}, maxNumberLimit int) {
	for {
		select {
			case <-quit:
				return
			default:
				numbersGen <- rand.Intn(maxNumberLimit)
		}
	}
}

func main() {
	var maxNumberLimit int
	var threadsCount int
	flag.IntVar(&maxNumberLimit, "n", 1, "Limit max number and count unique random numbers")
	flag.IntVar(&threadsCount, "threads", 1, "Count of threads")
	flag.Parse()

	if threadsCount < 1 {
		log.Fatalln("Flag -threads required greater 0")
	}

	if maxNumberLimit < 1 {
		return
	}

	rand.Seed(time.Now().UnixNano())
	numbersGen := make(chan int)
	quit := make(chan struct{})

	for i := 1; i <= threadsCount; i++ {
		go randomGenWorker(numbersGen, quit, maxNumberLimit)
	}

	go printer(numbersGen, maxNumberLimit, quit)
	<-quit
}
