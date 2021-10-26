package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
)

var DNServer = flag.String("dns", "1.1.1.1", "-dns <DNS IP> precise the DNS to use")
var nReq = flag.Int("nr", 1000, "-nr <number of requests> precise the total number of dns requests")
var nClient = flag.Int("nc", 10, "-nc <number of clients> precise the number of clients to dispatch queries on")

type ConcurrentIntSlice struct {
	sync.RWMutex
	items []float64
}

func (cs *ConcurrentIntSlice) Append(item float64) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
}

func main() {
	//Parse DNS, number of clients and number of requests
	flag.Parse()

	//Create resolver
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, *DNServer+":53")
		},
	}

	//Retrieve list of domains
	f, err := os.Open("data/top-1m-corrected.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	domainList := []string{}
	for scanner.Scan() {
		domainList = append(domainList, scanner.Text())
	}

	if *nReq > len(domainList) {
		fmt.Println("[CG-DNStester] Error: total number of requests must be < " + fmt.Sprint(len(domainList)))
	}

	//Launching clients
	var wg sync.WaitGroup
	size := int(*nReq / *nClient)
	var j int
	var timeList ConcurrentIntSlice
	fmt.Println("[CG-DNStester] Started requesting " + *DNServer)
	for i := 0; i < *nReq; i += size {
		j += size
		if j > *nReq {
			j = *nReq
		}
		wg.Add(1)
		go resolveAdresses(&timeList, &wg, *r, &domainList, i, j)
	}

	//Waiting for client to finish
	wg.Wait()
	fmt.Printf("[CG-DNStester] Total time = %v ms\n", timeList.items)
	mean, _ := stats.Mean(timeList.items)
	vari, _ := stats.StandardDeviation(timeList.items)
	min, _ := stats.Min(timeList.items)
	max, _ := stats.Max(timeList.items)
	fmt.Printf("[CG-DNStester] Max: %.1f ms, Min: %.1f ms, Mean: %.1f ms, StdDev: %.1f ms\n", max, min, mean, vari)
}

func resolveAdresses(timeList *ConcurrentIntSlice, wg *sync.WaitGroup, r net.Resolver, listOfDomains *[]string, i, j int) {
	defer wg.Done()
	t1 := time.Now()
	//Resolve requests
	for _, d := range (*listOfDomains)[i:j] {
		r.LookupHost(context.Background(), d)
	}
	ft := time.Now().Sub(t1).Milliseconds()
	timeList.Append(float64(ft))
}
