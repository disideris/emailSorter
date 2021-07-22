// Package concurrentcustomerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package concurrentcustomerimporter

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// Group type for holding chunks of strings
type Group []string

// Pool type for holding channels
type Pool struct {
	Work    chan string
	PreWork chan Group
	Result  chan string
}

// Domain type for holding domain names and its counts
type Domain struct {
	domainName string
	counts     int64
}

var domains []Domain
var line []byte
var bytesParts [][]byte
var bytesParts2 [][]byte
var domain string

// Functions for processing
func makePool(workers, groupFactor int, fn func(string) string) *Pool {
	var wg sync.WaitGroup
	p := &Pool{
		Work:    make(chan string, 100),
		PreWork: make(chan Group, 100),
		Result:  make(chan string, 100),
	}

	// launch workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			for g := range p.PreWork {
				for _, s := range g {
					p.Result <- fn(s)
				}
			}
			wg.Done()
		}()
	}

	// launch grouper
	wg.Add(1)
	go func() {
		var g Group
		for w := range p.Work {
			g = append(g, w)
			if len(g) == groupFactor {
				p.PreWork <- g
				g = nil
			}
		}
		p.PreWork <- g
		close(p.PreWork)
		wg.Done()
	}()

	// launch finisher
	go func() { wg.Wait(); close(p.Result) }()

	return p
}

// Function that extract the email domain from given file line
func extractDomain(line []byte) string {

	atIndex := bytes.IndexByte(line, '@')
	if atIndex == -1 {
		return ""
	}
	commaAfterAtIndex := bytes.IndexByte(line[atIndex+1:], ',') + atIndex + 1

	return string(line[atIndex+1 : commaAfterAtIndex])
}

// Function that reads the file and sends each line to the worker pool
func feedWorkers(filePath string, workPool chan string) {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//scanner.Scan()
	for scanner.Scan() {
		line := scanner.Bytes()
		workPool <- extractDomain(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	close(workPool)
}

// ImportCustomers function imports customers from file and returns a sorted array of email domains along with tis frequencies
func ImportCustomers(filePath string) []Domain {

	// defer profile.Start().Stop()

	// run in different cores in "parallel"
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	// Create a Pool instance
	pool := makePool(3, 10, func(s string) string {
		return s
	})

	// Read input file and send its lines to channel Work
	go feedWorkers(filePath, pool.Work)

	return makeSortedArray(countDomains(pool.Result))
}

func countDomains(resultPool chan string) map[string]int64 {

	// Map to hold the domains as keys and the domain counts as values
	domainMap := make(map[string]int64)

	// Put results to the hashmap for quick insertion O(1)
	for res := range resultPool {
		domainMap[res]++
	}
	return domainMap
}

// Function that converts a hashmap to a sorted slice of Domain structs
func makeSortedArray(domainMap map[string]int64) []Domain {

	//delete First line of file
	delete(domainMap, "")

	// convert hashmap to slice of structs
	for k, v := range domainMap {
		domains = append(domains, Domain{k, v})
	}

	// Use golang quicksort to sort slice with alphabetical order according to domain name
	sort.Slice(domains, func(i, j int) bool {
		switch strings.Compare(domains[i].domainName, domains[j].domainName) {
		case -1:
			return true
		case 1:
			return false
		}
		return domains[i].domainName > domains[j].domainName
	})
	return domains
}
