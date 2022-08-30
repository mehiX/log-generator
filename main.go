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
var generatorThroughput = 1000
var runFor = 1 * time.Second
var stats = true
var inputFile = ""

func init() {
	flag.StringVar(&inputFile, "in", "input.json", "Name of the input json file")
	flag.DurationVar(&runFor, "run", 1*time.Second, "Duration to run for")
	flag.IntVar(&generators, "g", 10, "Number of concurrent generators to use")
	flag.IntVar(&generatorThroughput, "gt", 1000, "How many entries per second should generate each generator")
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
		go generate(done, in, out, time.Second/time.Duration(generatorThroughput))
	}

	if stats {
		defer printDuration()(&counter)
	}

	for {
		select {
		case <-done.Done():
			return
		case s := <-out:
			counter++
			if !stats {
				fmt.Printf("%v\n", s)
			}
		}
	}

}

func printDuration() func(*int) {
	now := time.Now()

	return func(c *int) {
		fmt.Printf("Messages generated: %d [%v]\n", *c, time.Since(now))
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
