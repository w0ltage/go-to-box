package main

import (
  "bufio"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "strings"
)

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

func main() {
  hostsPath := "tmp/hosts"
  content, readResult := readFile(hostsPath)

  if readResult != nil {
    fmt.Printf("Error reading hosts: %v\n", readResult)
    return
  }

  fmt.Println(content)

  removeResult := removeDomain(hostsPath, content, "htb") 
  if removeResult != nil {
    fmt.Printf("Error removing domains from hosts: %v\n", readResult)
    return
  }

  fmt.Println("Domains removed successfully.")
}
