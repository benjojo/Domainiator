package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type LogPayload struct {
	Sucessful   bool
	Headers     http.Header
	DomainName  string
	RequestTime time.Duration
}

func worker(linkChan chan string, resultsChan chan LogPayload, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 45
	for url := range linkChan {
		start := time.Now()
		formattedurl := fmt.Sprintf("http://%s.com/", strings.TrimSpace(url))
		urlobj, e := http.Get(formattedurl)
		// fmt.Printf("BRB getting '%s'\n", formattedurl)
		if e == nil {
			elapsed := time.Since(start)

			Payload := LogPayload{
				DomainName:  strings.TrimSpace(url),
				Headers:     urlobj.Header,
				Sucessful:   true,
				RequestTime: elapsed,
			}
			resultsChan <- Payload
		} else {

			fakeheaders := make(http.Header)
			Payload := LogPayload{
				DomainName:  strings.TrimSpace(url),
				Headers:     fakeheaders,
				Sucessful:   false,
				RequestTime: 0,
			}

			resultsChan <- Payload
		}

	}

}

func Logger(resultChan chan LogPayload) {
	Database, e := GetDB()

	if e != nil {
		panic("Logger could not connect to the database")
	}

	for results := range resultChan {
		b, _ := json.Marshal(results)
		if results.Sucessful == true {
			Database.Exec("INSERT INTO `Domaniator`.`Results` (`Domain`, `Data`) VALUES (?, ?)", results.DomainName, string(b))
		} else {
			Database.Exec("INSERT INTO `Domaniator`.`Results` (`Domain`, `Data`) VALUES (?, ?)", results.DomainName, "f")
		}
	}
}

func main() {

	b, e := ioutil.ReadFile(os.Args[1])
	if e != nil {
		panic(e)
	}
	File := strings.Split(string(b), "\n")

	lCh := make(chan string)
	rCh := make(chan LogPayload)
	wg := new(sync.WaitGroup)
	go Logger(rCh)
	// Adding routines to workgroup and running then
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go worker(lCh, rCh, wg)
	}

	for _, link := range File {
		lCh <- link
	}
	// Closing channel (waiting in goroutines won't continue any more)
	close(lCh)
	wg.Wait()
}
