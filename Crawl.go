package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
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
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 15
	for url := range linkChan {
		start := time.Now()
		// Construct the HTTP request, I have to go this the rather complex way because I want
		// To add a useragent
		formattedurl := fmt.Sprintf("http://%s.com/", strings.TrimSpace(url))
		req, err := http.NewRequest("GET", formattedurl, nil)
		if err == nil {
			client := &http.Client{}
			client.CheckRedirect =
				func(req *http.Request, via []*http.Request) error {
					e := errors.New("can't go here because of golang bug")
					return e
				}
			req.Header.Set("User-Agent", "HTTP Header Survey By Benjojo (google benjojo) https://github.com/benjojo/Domainiator")
			urlobj, e := client.Do(req)
			// ioutil.ReadAll(urlobj.Body)
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

}

func Logger(resultChan chan LogPayload) {
	Database, e := GetDB()
	Query, _ := Database.Prepare("INSERT INTO `Domaniator`.`Results` (`Domain`, `Data`) VALUES (?, ?)")

	if e != nil {
		panic("Logger could not connect to the database")
	}

	for results := range resultChan {
		b, _ := json.Marshal(results)
		if results.Sucessful == true {
			Query.Exec(results.DomainName, string(b))
		} else {
			Query.Exec(results.DomainName, "f")
		}
	}
}

func main() {
	runtime.GOMAXPROCS(3)
	b, e := ioutil.ReadFile(os.Args[1])
	if e != nil {
		panic(e)
	}
	File := strings.Split(string(b), "\n")

	lCh := make(chan string)
	rCh := make(chan LogPayload, 100)
	wg := new(sync.WaitGroup)
	go Logger(rCh)
	// Adding routines to workgroup and running then
	for i := 0; i < 600; i++ {
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
