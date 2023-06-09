package main

import (
  "bufio"
  "fmt"
  "io/ioutil"
  "os"
  "strings"
  "flag"
)

func main() {
  var tld string
  var ip string
  var domain string
  var joinFlag bool
  // var domain string
  hostsPath := "tmp/hosts"
  content, readResult := readFile(hostsPath)

  // check if file exists or if user have permissions to read it
  if readResult != nil {
    fmt.Printf("Error reading hosts: %v\n", readResult)
    return
  } 

  flag.StringVar(&tld, "c", "", "Clear all hosts with a specific TLD")
  flag.BoolVar(&joinFlag, "j", false, "Join host to hosts file")
  flag.StringVar(&ip, "a", "", "Domain IP address")
  flag.StringVar(&domain, "d", "", "Domain name")
  flag.Parse()

  if tld != "" {
    removeResult := removeDomain(hostsPath, content, tld) 

    if removeResult != nil {
      fmt.Printf("Error removing domains from hosts: %v\n", removeResult)
      return
    } else {
      fmt.Println("\nDomains removed successfully.")
      return
    }
  }

  if joinFlag && (ip == "" || domain == "") {
    fmt.Printf("When using -j flag, both -a for IP-address and -d for domain name are required\n\n")
    flag.PrintDefaults()
    return
  } else if joinFlag && (ip != "" || domain != "") {
    addResult := addDomain(ip, domain, hostsPath)

    if addResult != nil {
      fmt.Printf("Error adding domains to hosts: %v\n", addResult)
      return
    } else {
      fmt.Println("\nDomain added successfully.")
      return 
    }
  }

  fmt.Printf("Nothing happened. Did you need something?\n\n")
  flag.PrintDefaults()
}

func readFile(filePath string) (string, error) {
  content, err := ioutil.ReadFile(filePath)
  return string(content), err
}

func removeDomain(hostsFile, content, substring string) error {
  // clear hosts file
  hostsOutput, err := os.Create(hostsFile)
  defer hostsOutput.Close()

  // iterate over each line in the content
  reader := strings.NewReader(content)
  scanner := bufio.NewScanner(reader)
  for scanner.Scan() {
    line := scanner.Text()

    // check if the line contains the substring
    if !strings.Contains(line, substring) {
      fmt.Fprintln(hostsOutput, line)
    }
  }

  if err := scanner.Err(); err != nil {
		return err
	}

  return err
}

func addDomain(ip, domain, hostsPath string) error {
  // open the hosts file in append mode
  file, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
  defer file.Close()

  if err != nil {
    fmt.Println("Error opening file:", err)
    return err
  }

  // write the new line to the end of the file
  line := fmt.Sprintf("%s    %s", ip, domain)
  _, err = fmt.Fprintln(file, line)

  if err != nil {
    fmt.Println("Error writing to file:", err)
    return err
  }

  return nil
}

