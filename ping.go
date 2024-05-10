package main

import (
	"fmt"
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

// runTcpPing runs tcp pings
func runTcpPing(port int, wgi, wgp *sync.WaitGroup, ipc chan netip.Addr, rec chan result) {
	go func() {
		for ip := range ipc {
			go func(ip netip.Addr) {
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), *timeout)
				if err == nil {
					defer conn.Close()
				}
				wgp.Add(1)
				rec <- result{ip: ip, err: err}
				wgi.Done()
			}(ip)
		}
	}()
}
