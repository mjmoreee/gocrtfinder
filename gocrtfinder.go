package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	targetURL := flag.String("u", "", "Target URL (google.com, microsoft.com)")
	flag.Parse()

	if *targetURL == "" {
		fmt.Println("[-] Target url doesn't exist, see --help for more info")
		fmt.Println("[+] Usage: ./ctrfinder -u target.com")
		os.Exit(1)
	}

	targetWord := strings.Split(*targetURL, ".")[0]
	fmt.Printf("[+] Getting the results from crt.sh with these keywords: *.%s.com / %s.*\n", targetWord, targetWord)
	fmt.Printf("[+] Start collecting info for: %s\n\n", targetWord)

	firstURL := fmt.Sprintf("https://crt.sh/?q=%%25.%s", *targetURL)
	secondURL := fmt.Sprintf("https://crt.sh/?q=%s.%%25", targetWord)

	firstResponse, err := http.Get(firstURL)
	if err != nil {
		fmt.Printf("Error fetching data from %s: %v\n", firstURL, err)
		os.Exit(1)
	}
	defer firstResponse.Body.Close()

	secondResponse, err := http.Get(secondURL)
	if err != nil {
		fmt.Printf("Error fetching data from %s: %v\n", secondURL, err)
		os.Exit(1)
	}
	defer secondResponse.Body.Close()

	firstContent, err := ioutil.ReadAll(firstResponse.Body)
	if err != nil {
		fmt.Printf("Error reading response body from %s: %v\n", firstURL, err)
		os.Exit(1)
	}

	secondContent, err := ioutil.ReadAll(secondResponse.Body)
	if err != nil {
		fmt.Printf("Error reading response body from %s: %v\n", secondURL, err)
		os.Exit(1)
	}

	firstSubdomains := extractSubdomains(string(firstContent))
	secondSubdomains := extractSubdomains(string(secondContent))

	allSubdomains := make(map[string]bool)
	for _, subdomain := range append(firstSubdomains, secondSubdomains...) {
		if !allSubdomains[subdomain] {
			fmt.Println(subdomain)
			allSubdomains[subdomain] = true
		}
	}

	outputFile := fmt.Sprintf("%s.txt", targetWord)
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", outputFile, err)
		os.Exit(1)
	}
	defer f.Close()

	for subdomain := range allSubdomains {
		if _, err := f.WriteString(subdomain + "\n"); err != nil {
			fmt.Printf("Error writing to file %s: %v\n", outputFile, err)
			os.Exit(1)
		}
	}

	if len(allSubdomains) > 0 {
		fmt.Printf("\n[+] Founded: %d domain related to %s\n", len(allSubdomains), targetWord)
		fmt.Printf("[+] Output file name: %s\n", outputFile)
	} else {
		fmt.Printf("\n[-] There's no subdomains related to: %s\n", targetWord)
	}
}

func extractSubdomains(content string) []string {
	var subdomains []string
	re := regexp.MustCompile(`[a-zA-Z0-9.-]+\.idn\.id`)
	matches := re.FindAllString(content, -1)
	for _, match := range matches {
		subdomains = append(subdomains, match)
	}
	return subdomains
}
