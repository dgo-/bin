package main

import "os"
import "bufio"
import "strings"
import "errors"
import "net"
import "time"
import "sync"
import "flag"
import "github.com/fatih/color"

// check if the server is reachable an the given port
func CheckHost(wg *sync.WaitGroup, server string, sec int) {
	timeout := time.Duration(sec)*time.Second
  defer wg.Done()

  c, e := net.DialTimeout("tcp", server, timeout)
	if e != nil {
    color.Red(server + " failed: " + e.Error())
	} else {
    color.Green(server + " working")
		c.Close()
	}
}

// check if the input is a valid address string
func CheckInput(i string) (string, error) {
  var err error
  s := strings.TrimSpace(i)

  if s == "" {
    err = errors.New(s + ": is empty")
  }
  if ! strings.Contains(s, ":") {
    err = errors.New(s + ": not contains a port")
  }
  return i, err
}

// read the servers from the given  file
func GetList(list string) ([]string, error) {
  var servers []string

  file, err := os.Open(list)
  if err == nil {
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      t, e := CheckInput(scanner.Text())
      if e == nil {
        servers = append(servers, t)
      }
    }
    err = scanner.Err()
    file.Close()
  }
  return servers, err
}

// loop all wait seconds read hosts file at the beginning of each run. 
func loop(wait int ,list string, timeout int) {
  var wg sync.WaitGroup
	for {
    servers, err := GetList(list)
    if err == nil {
      if len(servers) > 0 {
        for _, host := range servers {
          wg.Add(1)
          go CheckHost(&wg, host, timeout)
        }
        wg.Wait()
      } else {
        color.Yellow("no hosts in hosts file: " + list)
      }
    } else {
      color.Yellow("unable to read hosts file: " + err.Error())
    }
		time.Sleep(time.Duration(wait) * time.Second)
  }
}

// main function: parse commandline arguments and start main loop
func main() {
  wait    := flag.Int("wait", 10, "set the time until retry to connect")
  timeout := flag.Int("timeout", 1, "define the TCP timeout")
  list    := flag.String("list", os.Getenv("HOME") + "/hosts2check", "read hosts to check from file. Each line one host. Format: 127.0.0.1:80 www.google.com:443")
  flag.Parse()

	loop(*wait, *list, *timeout)
}
