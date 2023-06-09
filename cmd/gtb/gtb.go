package main

import (
    "bufio"
    "fmt"
    "flag"
    "io/ioutil"
    "log"
    "os"
    "strings"
)

func main() {
    // Initialize flags
    var tld string
    var ip string
    var domain string
    var joinFlag bool
    var replaceFlag bool

    // Define path to hosts file
    hostsPath := "/etc/hosts"

    // Read contents of hosts file
    content, readResult := readFile(hostsPath)
    if readResult != nil {
        fmt.Printf("Error reading hosts: %v\n", readResult)
        return
    }

    // Define command line flags
    flag.StringVar(&tld, "c", "", "Clear all hosts with a specified TLD")
    flag.BoolVar(&joinFlag, "j", false, "Join host to hosts file")
    flag.BoolVar(&replaceFlag, "r", false, "Remove all domains with specific TLD to replace them wtih another IP and domain")
    flag.StringVar(&ip, "a", "", "Domain IP address")
    flag.StringVar(&domain, "d", "", "Domain name")
    flag.Parse()

    // Handle clear flag
    if tld != "" && !replaceFlag {
        err := removeDomain(hostsPath, content, tld) 

        if err != nil {
            log.Fatalf("Error removing domains from hosts file: %v", err)
        }
        fmt.Println("\nDomains removed successfully.")
        return
    }

    // Handle join flag
    if joinFlag {
        if ip == "" || domain == "" {
            flag.Usage()
            log.Fatal("Both IP address (-a) and domain name (-d) are required when using -j flag")
        }

        err := addDomain(ip, domain, hostsPath)
        if err != nil {
            log.Default().Fatalf("Error adding domain to hosts file: %v", err)
        }

        fmt.Println("\nDomain added successfully")
        return
    }

    // Handle replace flag
    if replaceFlag {
        if ip == "" || domain == "" || tld == "" {
            flag.Usage()
            log.Fatal("IP address (-a), domain name (-d) and TLD (-c) are required when using -r flag")
        }

        err := removeDomain(hostsPath, content, tld)
        if err != nil {
            log.Fatalf("Error removing domains from hosts file: %v", err)
        }

        err = addDomain(hostsPath, ip, domain)
        if err != nil {
            log.Fatalf("Error adding domain to hosts file: %v", err)
        }

        fmt.Println("\nDomain replaced successfully")
        return
    }

    // If no flag is specified, print usage information
    fmt.Printf("You don't need to go to box?\n\n")
    flag.Usage()
}

// readFile reads the contents of a file and returns it as a string
func readFile(filePath string) (string, error) {
    content, err := ioutil.ReadFile(filePath)
    return string(content), err
}

// removeDomain removes all lines from a file that contain a specified substring
func removeDomain(hostsFile, content, substring string) error {
    // Open hosts file for writing
    hostsOutput, err := os.Create(hostsFile)
    if err != nil {
        return err
    }
    defer hostsOutput.Close()

    // Iterate over each line in the content
    reader := strings.NewReader(content)
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        line := scanner.Text()

        // Check if the line contains the substring
        if !strings.Contains(line, substring) {
            fmt.Fprintln(hostsOutput, line)
        }
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    return nil
}

// addDomain appends a new line to a file with a specified IP address and domain name
func addDomain(hostsPath, ip, domain string) error {
    // Open hosts file for appending
    file, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write new line to end of the file
    line := fmt.Sprintf("%s    %s", ip, domain)
    _, err = fmt.Fprintln(file, line)

    if err != nil {
        return err
    }

    return nil
}

