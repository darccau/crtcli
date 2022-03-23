package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type CrtDomains struct {
	NameValue string `json:"name_value"`
}

func getRawDomains(baseDomain string) []CrtDomains {
	var subdomains []CrtDomains

	url := "https://crt.sh/?q=" + baseDomain + "&output=json"

	response, error := http.Get(url)

	if error != nil {
		log.Fatal(error)
	}
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	if error != nil {
		log.Fatal(error)
	}

	json.Unmarshal(body, &subdomains)

	return subdomains
}

func cleanSubdomains(rawDomains []CrtDomains) map[string]bool {
	uniqueDomains := make(map[string]bool)

	for _, subdomain := range rawDomains {
		found, _ := regexp.MatchString("@|\\*", subdomain.NameValue)
		if !found {
			uniqueDomains[subdomain.NameValue] = true
		}
	}

	return uniqueDomains
}

func displayDomains(cleanSubdomains map[string]bool) {
	for subdomains, _ := range cleanSubdomains {
		fmt.Println(subdomains)
	}
}

func main() {
	baseDomain := os.Args[1]
	rawSubdomains := getRawDomains(baseDomain)
	subdomains := cleanSubdomains(rawSubdomains)
	displayDomains(subdomains)
}
