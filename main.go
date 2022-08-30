package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var generators = 10
var delay = 1 * time.Millisecond
var runFor = 1 * time.Second
var stats = true
var inputFile = ""

func init() {
	flag.StringVar(&inputFile, "in", "input.json", "Name of the input json file")
	flag.DurationVar(&runFor, "run", 1*time.Second, "Duration to run for")
	flag.IntVar(&generators, "g", 10, "Number of concurrent generators to use")
	flag.BoolVar(&stats, "stats", false, "False to print the actual messages, True to print only stats")

	flag.Parse()
}

func main() {

	if inputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	b, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	messages := make([]interface{}, 0)
	if err := json.Unmarshal(b, &messages); err != nil {
		panic(err)
	}

	out := make(chan interface{}, generators)
	done, cancel := context.WithTimeout(context.Background(), runFor)
	defer cancel()

	counter := 0

	in := messageGenerator(done, messages)

	for i := 0; i < generators; i++ {
		go generate(done, in, out, delay)
	}

	if stats {
		defer printDuration()()
	}

	for {
		select {
		case <-done.Done():
			if stats {
				fmt.Printf("Counter: %d\n", counter)
			}
			return
		case s := <-out:
			counter++
			if !stats {
				fmt.Printf("%v\n", s)
			}
		}
	}

}

func printDuration() func() {
	now := time.Now()

	return func() {
		fmt.Printf("Duration: %v\n", time.Since(now))
	}
}

func messageGenerator(ctx context.Context, m []interface{}) <-chan interface{} {
	c := make(chan interface{})
	rand.Seed(time.Now().UnixMilli())

	go func() {
		defer close(c)

		l := len(m)
		for {
			select {
			case <-ctx.Done():
				return
			case c <- m[rand.Intn(l)]:
			}
		}
	}()

	return c
}

func generate(done context.Context, in <-chan interface{}, out chan<- interface{}, d time.Duration) {
	tckr := time.NewTicker(d)
	defer tckr.Stop()

	for {
		select {
		case <-done.Done():
			return
		case <-tckr.C:
			out <- <-in
		}
	}
}
