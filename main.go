package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var client *http.Client

var concurrency, timeout int

func main() {

	flag.IntVar(&concurrency, "c", 50, "concurrency")
	flag.IntVar(&timeout, "t", 5, "timeout in seconds")
	flag.Parse()

	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        concurrency,
			MaxIdleConnsPerHost: concurrency,
			MaxConnsPerHost:     concurrency,
			DialContext: (&net.Dialer{
				Timeout: time.Duration(time.Duration(timeout) * time.Second),
			}).DialContext,
			TLSHandshakeTimeout: time.Duration(time.Duration(timeout) * time.Second),
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS10,
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Duration(time.Duration(timeout) * time.Second),
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	inputURLS := []string{}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		input := sc.Text()
		inputURLS = append(inputURLS, input)
	}

	var wg sync.WaitGroup
	urlsChan := make(chan string)
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for uri := range urlsChan {

				urlToCheck, redirectSubdomain := createRedirectURL(uri)

				// request and check for 307
				redirectDidIndeedHappen, redirectURL := redirectHappened(urlToCheck, redirectSubdomain)
				if redirectDidIndeedHappen {
					fmt.Printf("%s\t[Location: %s]\n", urlToCheck, redirectURL)
				}
			}
			wg.Done()
		}()
	}

	for _, uri := range inputURLS {
		urlsChan <- uri
	}

	close(urlsChan)
	wg.Wait()
}
