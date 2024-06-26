package main

import (
	"flag"
	"fmt"
	"net/netip"
	"os"
	"sync"
	"time"

	"github.com/digineo/go-ping"
	"github.com/fatih/color"
)

var timeout = flag.Duration("timeout", time.Second, "ping timeout")

func main() {
	flag.Parse()

	nets := flag.Args()
	if len(nets) == 0 {
		fmt.Fprintln(os.Stderr, "at least one ip, or ip-ip range, or network/ip is required")
		os.Exit(1)
	}

	ips := map[netip.Addr]struct{}{}

	for _, n := range nets {
		ip, err := netip.ParseAddr(n)
		if err == nil {
			ips[ip] = struct{}{}
			continue
		}

		ipr, err := parseIPRange(n)
		if err == nil {
			for ip := ipr.from; ip.Compare(ipr.to) <= 0; ip = ip.Next() {
				ips[ip] = struct{}{}
			}
			continue
		}

		ipp, err := netip.ParsePrefix(n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "network/ip %s is not a valid ip, ip range or network\n", n)
			os.Exit(1)
		}

		ipp = ipp.Masked()
		for ip := ipp.Addr(); ipp.Contains(ip); ip = ip.Next() {
			ips[ip] = struct{}{}
		}
	}

	p, err := ping.New("0.0.0.0", "::")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ipc := make(chan netip.Addr, 10)
	rec := make(chan result, 10)
	wgi, wgp := sync.WaitGroup{}, sync.WaitGroup{}

	runIcmpPing(p, &wgi, &wgp, ipc, rec)

	for ip := range ips {
		ipc <- ip
		wgi.Add(1)
	}

	go func() {
		for res := range rec {
			fmt.Print(res.ip, "\t")
			if res.err == nil {
				color.Green("UP")
			} else {
				color.Red("DOWN")
			}
			wgp.Done()
		}
	}()

	wgi.Wait()
	close(ipc)
	wgp.Wait()
	close(rec)
}
