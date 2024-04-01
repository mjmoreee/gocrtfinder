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
    targetURL := flag.String("u", "", "Target URL (e.g., google.com, microsoft.com)")
    outputDir := flag.String("o", ".", "Output directory path")
    flag.Parse()

    if *targetURL == "" {
        fmt.Println("[-] Target URL is not provided, see --help for more information")
        fmt.Println("[+] Usage: ./ctrfinder -u target.com -o /path/to/output")
        os.Exit(1)
    }

    targetWord := strings.Split(*targetURL, ".")[0]
    fmt.Printf("[+] Getting the results from crt.sh with these keywords: *.%s / %s.*\n", targetWord, targetWord)
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

    firstSubdomains := extractSubdomains(string(firstContent), *targetURL)
    secondSubdomains := extractSubdomains(string(secondContent), targetWord)

    allSubdomains := make(map[string]bool)
    for _, subdomain := range append(firstSubdomains, secondSubdomains...) {
        if !allSubdomains[subdomain] {
            allSubdomains[subdomain] = true
        }
    }

    outputFile := fmt.Sprintf("%s/%s.txt", *outputDir, targetWord)
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
        fmt.Printf("\n[+] Founded: %d domains related to %s\n", len(allSubdomains), targetWord)
        fmt.Printf("[+] Output file name: %s\n", outputFile)
    } else {
        fmt.Printf("\n[-] There are no subdomains related to: %s\n", targetWord)
    }
}

func extractSubdomains(content string, targetDomain string) []string {
    var subdomains []string
    re := regexp.MustCompile(`[a-zA-Z0-9.-]+\.` + regexp.QuoteMeta(targetDomain))
    matches := re.FindAllString(content, -1)
    for _, match := range matches {
        subdomains = append(subdomains, match)
    }
    return subdomains
}
