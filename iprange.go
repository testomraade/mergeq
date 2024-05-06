package main

import (
	"fmt"
	"net/netip"
	"strings"
)

// ipRange holds two ip addr
type ipRange struct {
	from netip.Addr
	to   netip.Addr
}

// parseIPRange parses two ip addres seperated by a hyphen
func parseIPRange(s string) (ipRange, error) {
	var r ipRange
	h := strings.IndexByte(s, '-')
	if h == -1 {
		return r, fmt.Errorf("no hyphen in range %q", s)
	}
	from, to := s[:h], s[h+1:]
	var err error
	r.from, err = netip.ParseAddr(from)
	if err != nil {
		return r, fmt.Errorf("invalid From IP %q in range %q", from, s)
	}
	r.to, err = netip.ParseAddr(to)
	if err != nil {
		return r, fmt.Errorf("invalid To IP %q in range %q", to, s)
	}
	if !r.isValid() {
		fmt.Println(r.from.IsValid(), r.to.IsValid(), r.from.Compare(r.to), r.to.Less(r.from))
		return r, fmt.Errorf("range %v to %v not valid", r.from, r.to)
	}
	return r, nil
}

func (r ipRange) isValid() bool {
	return r.from.IsValid() &&
		r.to.IsValid() &&
		!r.to.Less(r.from)
}
