package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
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
	DNSIP       string
	DomainName  string
	RequestTime time.Duration
	StatusCode  int
}

func worker(linkChan chan string, resultsChan chan LogPayload, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 15
	http.DefaultTransport.(*http.Transport).DisableKeepAlives = true

	for url := range linkChan {
		start := time.Now()
		// Construct the HTTP request, I have to go this the rather complex way because I want
		// To add a useragent
		tld := ""
		if *presumecom {
			tld = ".com"
		}

		formattedurl := fmt.Sprintf("http://%s%s%s", strings.TrimSpace(url), tld, *pathtoquery)
		req, err := http.NewRequest("GET", formattedurl, nil)

		if err != nil {
			fmt.Printf("Issue with URL provided, %s makes no sence to net/url\n", formattedurl)
			continue
		}

		client := &http.Client{}
		client.CheckRedirect =
			func(req *http.Request, via []*http.Request) error {
				e := errors.New("can't go here because of golang bug")
				return e
			}
		req.Header.Set("User-Agent", *useragent)

		// Avoid calling our own loopback, or calling on anything that does not have
		// DNS responce.
		ip, err := net.LookupIP(fmt.Sprintf("%s%s", strings.TrimSpace(url), tld))
		if err != nil || len(ip) < 1 || strings.HasPrefix("127.", ip[0].String()) || strings.HasPrefix("0.", ip[0].String()) {
			continue
		}
		urlobj, e := client.Do(req)
		if e == nil {
			elapsed := time.Since(start)
			if *saveoutput && urlobj.StatusCode == 200 {
				b, e := ioutil.ReadAll(urlobj.Body)
				if e == nil {
					os.Mkdir(fmt.Sprintf("./%s", strings.TrimSpace(url)[0]), 744)
					ioutil.WriteFile(fmt.Sprintf("./%s/%s.%s", strings.TrimSpace(url)[0], strings.TrimSpace(url), *pathtoquery), b, 744)
				}
			} else {
				urlobj.Body.Close()
			}

			Payload := LogPayload{
				DomainName:  strings.TrimSpace(url),
				Headers:     urlobj.Header,
				Sucessful:   true,
				DNSIP:       ip[0].String(),
				RequestTime: elapsed,
				StatusCode:  urlobj.StatusCode,
			}
			resultsChan <- Payload
		} else {

			fakeheaders := make(http.Header)
			Payload := LogPayload{
				DomainName:  strings.TrimSpace(url),
				Headers:     fakeheaders,
				Sucessful:   false,
				RequestTime: 0,
				StatusCode:  urlobj.StatusCode,
			}

			resultsChan <- Payload
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

var pathtoquery *string
var saveoutput *bool
var presumecom *bool
var databasestring *string
var useragent *string

func main() {
	runtime.GOMAXPROCS(3)
	inputfile := flag.String("input", "", "The file that will be read.")
	pathtoquery = flag.String("querypath", "/", "The path that will be queried.")
	saveoutput = flag.Bool("savepage", false, "Save the file that is downloaded to disk")
	presumecom = flag.Bool("presumecom", true, "Presume that the file lines need .com adding to them")
	concurrencycount := flag.Int("concount", 600, "How many go routines you want to start")
	databasestring = flag.String("dbstring", "root:@tcp(127.0.0.1:3306)/Domaniator", "What to connect to the database with")
	useragent = flag.String("ua", "HTTP Header Survey By Benjojo +https://github.com/benjojo/Domainiator", "What UA to send the request with")

	flag.Parse()

	if *inputfile == "" {
		fmt.Println("No input file, put one in with -input")
		os.Exit(0)
	}

	b, e := ioutil.ReadFile(*inputfile)
	if e != nil {
		panic(e)
	}
	File := strings.Split(string(b), "\n")

	lCh := make(chan string)
	rCh := make(chan LogPayload, 100)
	wg := new(sync.WaitGroup)
	go Logger(rCh)
	// Adding routines to workgroup and running then
	for i := 0; i < *concurrencycount; i++ {
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
