package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

const (
	hostsFile = "hosts.txt"
)

// IPStatus holds the address and the number of periods since the last
// successful ping attempt on the address
type IPStatus struct {
	IP     string
	Status int
}

// getHosts gets host addresses split by newlines from hosts.txt
func getHosts(file string) []IPStatus {
	out, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("%s: couldn't open up %s\n", err, file)
	}
	input := string(out)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	ips := make([]IPStatus, len(lines))
	for i := range lines {
		ips[i] = IPStatus{lines[i], 0}
	}
	return ips
}

// ping checks to see if a host is responsive
func (ip *IPStatus) ping() {
	err := exec.Command("ping", "-w 4", "-c 1", ip.IP).Run()
	if err != nil {
		ip.Status++
	} else {
		ip.Status = 0
	}
}

// handleOutage is called when a host is not responding for 2
// consecutive periods
func handleOutage(ip string) {
	fmt.Println("Outage handled!", ip)
}
func main() {
	ips := getHosts(hostsFile)
	for {
		for i := range ips {
			ips[i].ping()
			fmt.Printf("%s: %v\n", ips[i].IP, ips[i].Status)
			if ips[i].Status == 2 {
				handleOutage(ips[i].IP)
			}
		}

		fmt.Println(time.Now().Format(time.RubyDate))
		time.Sleep(time.Second * 5)
	}
}
