// Package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/pkg/profile"
)

//Domain type for holding domain name and its counts
type Domain struct {
	domainName string
	counts     int64
}

var domains []Domain
var line []byte

// Function that extract the email domain from given file line
func extractDomain(line []byte) string {

	atIndex := bytes.IndexByte(line, '@')
	if atIndex == -1 {
		return ""
	}
	commaAfterAtIndex := bytes.IndexByte(line[atIndex+1:], ',') + atIndex + 1

	return string(line[atIndex+1 : commaAfterAtIndex])
}

// Function that counts the number of email domains from the whole file
func countDomains(file io.Reader) map[string]int64 {

	domainMap := make(map[string]int64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		domainMap[extractDomain(line)]++
	}
	if e := scanner.Err(); e != nil {
		log.Fatal(e)
	}
	delete(domainMap, "")

	return domainMap
}

// ImportCustomers function imports customers from file and returns a sorted array of email domains along with tis frequencies
func ImportCustomers(filePath string) []Domain {

	defer profile.Start().Stop()

	file, e := os.Open(filePath)
	if e != nil {
		log.Fatal(e)
	}
	defer file.Close()

	return makeSortedArray(countDomains(file))
}

// Function that converts a hashmap to a sorted slice of Domain structs
func makeSortedArray(domainMap map[string]int64) []Domain {

	for k, v := range domainMap {
		domains = append(domains, Domain{k, v})
	}
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
