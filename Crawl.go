package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	// "time"
)

func worker(linkChan chan string, resultsChan chan string, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()

	for url := range linkChan {
		formattedurl := fmt.Sprintf("http://%s.com/", strings.TrimSpace(url))
		urlobj, e := http.Get(formattedurl)
		// fmt.Printf("BRB getting '%s'\n", formattedurl)
		if e == nil {
			b, _ := json.Marshal(urlobj.Header)
			// fmt.Println(string(b))
			resultsChan <- string(b)
		}
	}

}

func Logger(resultChan chan string) {
	Database, e := GetDB()
	if e != nil {
		panic "Logger could not connect to the database"
	}
	for results := range resultChan {
		fmt.Printf("BOOM %s", results)
	}
}

func main() {
	b, e := ioutil.ReadFile("./list.txt")
	if e != nil {
		panic(e)
	}
	File := strings.Split(string(b), "\n")

	lCh := make(chan string)
	rCh := make(chan string)
	wg := new(sync.WaitGroup)
	go Logger(rCh)
	// Adding routines to workgroup and running then
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go worker(lCh, rCh, wg)
	}
	var results string
	for _, link := range File {
		lCh <- link
	}
	// Closing channel (waiting in goroutines won't continue any more)
	close(lCh)
	wg.Wait()
}
