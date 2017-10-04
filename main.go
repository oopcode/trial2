package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	maxWorkers = 10
	logger     = log.New(os.Stdout, "", 0)
)

type (
	Input  chan string
	Output chan int
)

func main() {
	var (
		totalCount = 0
		input      = getInput()
		output     = getOutput(input, maxWorkers)
	)

	for pageCount := range output {
		totalCount += pageCount
	}

	logger.Printf("%d\n", totalCount)
}

// Creates and returns an input channel. Sets up a goroutine that reads urls
// from standard input and sends them to the input channel. Urls should be
// separated by newlines.
func getInput() Input {
	var (
		input  = make(Input)
		reader = bufio.NewReader(os.Stdin)
	)

	go func() {
		for {
			url, _, err := reader.ReadLine()
			if err != nil {
				close(input)
				break
			}

			input <- string(url)
		}
	}()

	return input
}

// getOutput sets up a channel or results and starts spawning workers that
// write to it. At most `maxWorkers` can be spawned at the same time.
func getOutput(input Input, concurrency int) Output {
	out := make(Output)

	go func() {
		// This semaphore is used to control the number of workers
		// and to wait for all of them to finish.
		sph := make(chan struct{}, concurrency)

		for val := range input {
			sph <- struct{}{}
			go func(val string) {
				out <- processOne(val)
				<-sph
			}(val)
		}

		// Stops when all jobs are done.
		for i := 0; i < concurrency; i++ {
			sph <- struct{}{}
		}

		close(out)
	}()

	return out
}

// processOne loads a web page and counts occurrences of the word "Go".
func processOne(url string) int {
	var (
		buf       = bytes.NewBuffer([]byte{})
		resp, err = http.Get(string(url))
	)
	if err != nil {
		logger.Printf("Failed to get %s: %s\n", url, err)
		return 0
	}
	defer resp.Body.Close()

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		logger.Printf("Failed to get %s: %s\n", url, err)
		return 0
	}

	var (
		page      = buf.String()
		pageCount = strings.Count(page, "Go")
	)

	logger.Printf("%s: %d", url, pageCount)

	return pageCount
}
