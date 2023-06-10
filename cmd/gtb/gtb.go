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
    var tldFlag string
    var ip string
    var domain string
    var appendFlag bool
    var replaceFlag bool

    // Define path to hosts file
    hostsPath := "tmp/hosts"

    // Read contents of hosts file
    content, readResult := readFile(hostsPath)
    if readResult != nil {
        fmt.Printf("Error reading hosts: %v\n", readResult)
        return
    }

    // Define command line flags
    flag.StringVar(&tldFlag, "rm", "", "Mode to remove all domains with a specified TLD")
    flag.BoolVar(&appendFlag, "add", false, "Mode to add host to hosts file")
    flag.BoolVar(&replaceFlag, "re", false, "Mode to remove all domains with specific TLD to replace them wtih another IP and domain")
    flag.StringVar(&ip, "i", "", "Value for domain IP address")
    flag.StringVar(&domain, "d", "", "Value for domain name")

    // Custom usage
    flag.Usage = func() {
        flagSet := flag.CommandLine
        fmt.Printf(
            "\nAdd, remove or replace hosts in the hosts file" +
            "\nUsage: %s { mode } { argument(s) } \n",
            "gtb")
        order := []string{"rm", "add", "re", "i", "d"}
        for _, name := range order {
            flag := flagSet.Lookup(name)
            fmt.Printf("-%s\t%s\n", flag.Name,flag.Usage)
        }
        fmt.Printf(
            "\nThere are only 3 types of program execution scenarios:\n" +
            "%[1]s -rm <tld>                                Remove all <tld> domains from hosts file\n" +
            "%[1]s -add -i <IP> -d <domain>                 Add <IP> address and <domain>\n" +
            "%[1]s -re -i -rm <tld> -i <IP> -d <domain>     Remove all <tld> domains and add <IP> with <domain>\n\n",
            "gtb")
    }

    flag.Parse()

    // Handle -rm (remove) flag
    if tldFlag != "" && !replaceFlag {
        err := removeDomain(hostsPath, content, tldFlag) 

        if err != nil {
            log.Printf("Error removing domains from hosts file: %v", err)
            flag.Usage()
            os.Exit(1)
        }
        fmt.Println("Domains removed successfully.")
        return
    }

    // Handle -add flag
    if appendFlag {
        if ip == "" || domain == "" {
            log.Println("Both IP address (-i) and domain name (-d) are required when using -add flag")
            flag.Usage()
            os.Exit(1)
        }

        err := addDomain(hostsPath, ip, domain)
        if err != nil {
            log.Fatalf("Error adding domain to hosts file: %v", err)
        }

        fmt.Println("Domain added successfully")
        return
    }

    // Handle -re (replace) flag
    if replaceFlag {
        if ip == "" || domain == "" || tldFlag == "" {
            log.Println("IP address (-i), domain name (-d) and TLD (-rm) are required when using -re flag")
            flag.Usage()
            os.Exit(1)
        }

        err := removeDomain(hostsPath, content, tldFlag)
        if err != nil {
            log.Fatalf("Error removing domains from hosts file: %v", err)
        }

        err = addDomain(hostsPath, ip, domain)
        if err != nil {
            log.Fatalf("Error adding domain to hosts file: %v", err)
        }

        fmt.Println("Domain replaced successfully")
        return
    }

    // If no flag is specified, print usage information
    fmt.Printf("I don't understand you.\nDon't you need to get into the box?\n\n")
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

