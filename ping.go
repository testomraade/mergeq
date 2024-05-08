package main

import (
	"net"
	"net/netip"
	"sync"

	"github.com/digineo/go-ping"
)

// result represents run ping error and endpoint ip
type result struct {
	ip  netip.Addr
	err error
}

// runIcmpPing runs icmp pings
func runIcmpPing(p *ping.Pinger, wgi, wgp *sync.WaitGroup, ipc chan netip.Addr, rec chan result) {
	go func() {
		for ip := range ipc {
			go func(ip netip.Addr) {
				_, err := p.Ping(&net.IPAddr{IP: ip.AsSlice()}, *timeout)
				wgp.Add(1)
				rec <- result{ip: ip, err: err}
				wgi.Done()
			}(ip)
		}
	}()
}
