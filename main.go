package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/digineo/go-ping"
	"github.com/fatih/color"
	"inet.af/netaddr"
)

var timeout = flag.Duration("timeout", time.Second, "ping timeout")

func main() {
	flag.Parse()

	nets := flag.Args()
	if len(nets) == 0 {
		fmt.Fprintln(os.Stderr, "at least one network/ip is required")
		os.Exit(1)
	}

	ips := map[netaddr.IP]struct{}{}

	for _, n := range nets {
		ip, err := netaddr.ParseIP(n)
		if err == nil {
			ips[ip] = struct{}{}
			continue
		}

		ipr, err := netaddr.ParseIPRange(n)
		if err == nil {
			for ip := ipr.From(); ip.Compare(ipr.To()) <= 0; ip = ip.Next() {
				ips[ip] = struct{}{}
			}
			continue
		}

		ipp, err := netaddr.ParseIPPrefix(n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "network/ip %s is not a valid ip, ip range or network\n", n)
			os.Exit(1)
		}

		ipr = ipp.Range()
		for ip := ipr.From(); ip.Compare(ipr.To()) <= 0; ip = ip.Next() {
			ips[ip] = struct{}{}
		}
	}

	p, err := ping.New("0.0.0.0", "::")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ipc := make(chan netaddr.IP, 10)
	wg := sync.WaitGroup{}
	writeLock := sync.Mutex{}

	go func() {
		for ip := range ipc {
			go func(ip netaddr.IP) {
				// fmt.Println(ip)
				_, err := p.Ping(ip.IPAddr(), *timeout)

				writeLock.Lock()
				fmt.Print(ip, "\t")
				if err == nil {
					color.Green("UP")
				} else {
					color.Red("DOWN")
				}
				writeLock.Unlock()

				wg.Done()
			}(ip)
		}
	}()

	for ip := range ips {
		ipc <- ip
		wg.Add(1)
	}

	wg.Wait()
}
