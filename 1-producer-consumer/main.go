//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(wg *sync.WaitGroup, stream Stream, tweets chan *Tweet) {
	defer wg.Done()

	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return
		}

		tweets <- tweet
	}
}

func consumer(wg *sync.WaitGroup, tweets <-chan *Tweet) {
	defer wg.Done()

	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {

	// This also works without making consumer a go func
	// but I thought I'd practice with wait groups.
	// Both functions are running concurrently, however,
	// consumer only finishes after producer does (as it's dependant on it),
	// making the wait at the end a bit unnecessary.
	wg := new(sync.WaitGroup)

	// Two go funcs hence the number two here and two "defer wg.Done()"
	wg.Add(2)

	start := time.Now()
	stream := GetMockStream()

	// Producer
	tweets := make(chan *Tweet)
	go producer(wg, stream, tweets)

	// Consumer
	go consumer(wg, tweets)

	// wait for the go routines to finish
	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
