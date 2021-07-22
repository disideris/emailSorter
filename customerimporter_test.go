package customerimporter

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestExtractDomainFromLine(t *testing.T) {

	t.Log("Testing extractDomainFromLine function...")

	line := []byte("Norma,Allen,nallen8@cnet.com,Female,168.67.162.1")
	expectedDomain := "cnet.com"
	actualDomain := extractDomain(line)

	if expectedDomain != actualDomain {
		t.Error("Domain is not extracted correctly")
	}
}

func TestMakeSortedArray(t *testing.T) {

	t.Log("Testing makeSortedArray function...")
	domainMap := make(map[string]int64)
	domainMap["zdnet.com"] = 5
	domainMap["about.com"] = 13
	domainMap["google.com"] = 7

	actualDomains := makeSortedArray(domainMap)
	expectedDomains := []Domain{{"about.com", 13}, {"google.com", 13}, {"zdnet.com", 5}}

	flag := true

	for i := 0; i < len(actualDomains); i++ {
		if actualDomains[i].domainName != expectedDomains[i].domainName {
			flag = false
		}
	}
	if !flag {
		t.Error("makeSortedArray function does not sort map")
	}
}

func TestCountDomains(t *testing.T) {

	t.Log("Testing testCountDomains function...")

	file, err := os.Create("/tmp/testcustomers.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
		return
	}

	fmt.Fprintf(file, "first_name,last_name,email,gender,ip_address\n")
	fmt.Fprintf(file, "Deborah,Moreno,dmorenohn@yahoo.com,Female,198.222.187.18\n")
	fmt.Fprintf(file, "Christina,Vasquez,cvasquezl0@zdnet.com,Female,91.236.117.57\n")
	fmt.Fprintf(file, "Phyllis,Lawrence,plawrence7x@yahoo.com,Female,48.210.189.102\n")
	fmt.Fprintf(file, "Willie,Ford,wford96@yahoo.com,Male,224.244.143.184\n")
	fmt.Fprintf(file, "Jane,Cunningham,jcunninghamcu@zdnet.com,Female,127.9.53.198\n")
	fmt.Fprintf(file, "Jonathan,Meyer,jmeyerj2@yahoo.com,Male,202.58.253.34\n")
	fmt.Fprintf(file, "Randy,Nichols,rnicholsdi@zdnet.com,Male,234.51.144.226\n")
	fmt.Fprintf(file, "Sandra,Gilbert,sgilbertk3@yahoo.com,Female,83.6.172.127\n")

	file.Close()

	file, e := os.Open("/tmp/testcustomers.csv")
	if e != nil {
		log.Fatal(e)
	}
	defer file.Close()

	flag := true

	expectedDomainMap := make(map[string]int64)
	expectedDomainMap["yahoo.com"] = 5
	expectedDomainMap["zdnet.com"] = 3

	actualDomainMap := countDomains(file)

	for k, v := range expectedDomainMap {
		if v != actualDomainMap[k] {
			flag = false
		}
	}

	if !flag {
		t.Error("Domain count is not as expected")
	}

	if actualDomainMap[""] != 0 {
		t.Error("Empty string is found in map")
	}

	// err2 := os.Remove("/tmp/testcustomers.csv")
	// if err2 != nil {
	// 	log.Fatal(e)
	// }
}
