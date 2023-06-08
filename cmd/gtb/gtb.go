package main

import (
  "bufio"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "strings"
  "flag"
)

func main() {
  hostsPath := "tmp/hosts"
  content, readResult := readFile(hostsPath)

  // check if file exists or if user have permissions to read it
  if readResult != nil {
    fmt.Printf("Error reading hosts: %v\n", readResult)
    return
  } else {
    fmt.Println(content)
  }

  var tld string
  flag.StringVar(&tld, "c", "", "Clear all hosts with a specific TLD")
  flag.Parse()
  
  if tld != "" {
    removeResult := removeDomain(hostsPath, content, tld) 

    if removeResult != nil {
      fmt.Printf("Error removing domains from hosts: %v\n", removeResult)
      return
    } else {
      fmt.Println("\nDomains removed successfully.")
    }
  }

}

func readFile(filePath string) (string, error) {
  content, err := ioutil.ReadFile(filePath)

  if err != nil {
      log.Fatal(err)
  }

  return string(content), nil
}

func removeDomain(hostsFile, content, substring string) error {

  // clear hosts file
  hostsOutput, err := os.Create(hostsFile)
  if err != nil {
    log.Fatal(err)
  }

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

  return nil
}

