package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func scanDir(dirPath string, semaphore chan struct{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out) // Close only when ALL files are processed

		dir, err := os.Open(dirPath)
		if err != nil {
			fmt.Println("Error opening directory:", err)
			return
		}
		defer dir.Close()

		for {
			semaphore <- struct{}{}
			entries, err := dir.ReadDir(1)
			if err != nil {
				<-semaphore // Release semaphore
				if err.Error() != "EOF" {
					fmt.Println("Error reading directory:", err)
				}
				break
			}

			// Send filename then release semaphore
			if len(entries) > 0 {
				out <- entries[0].Name()
			}
		}
	}()
	return out
}

func main() {
	dir := "/data/movies"
	semaphore := make(chan struct{}, 10)
	fileChannel := scanDir(dir, semaphore)
	var wg sync.WaitGroup

	for fileName := range fileChannel {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// time.Sleep(3 * time.Second)
			fmt.Println("Processing " + name)
			// time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
			time.Sleep(1 * time.Second)
			fmt.Println("Processing " + name + " is completed")
		}(fileName)
	}
	wg.Wait()
	fmt.Println("All files processed!")

}
