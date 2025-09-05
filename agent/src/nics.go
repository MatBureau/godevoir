package main

import (
	"time"

	"github.com/shirou/gopsutil/v4/net"
)

type NicRate struct {
	Name   string   `json:"name"`
	MTU    int      `json:"mtu"`
	Addr   []string `json:"addr"`   // IPv4/IPv6
	RxBps  float64  `json:"rx_bps"` // octets/s
	TxBps  float64  `json:"tx_bps"`
	RxMbps float64  `json:"rx_mbps"`
	TxMbps float64  `json:"tx_mbps"`
	Up     bool     `json:"up"`
}

func isUp(flags []string) bool {
	for _, f := range flags {
		if f == "up" || f == "UP" {
			return true
		}
	}
	return false
}

func isLoopback(flags []string) bool {
	for _, f := range flags {
		if f == "loopback" || f == "LOOPBACK" {
			return true
		}
	}
	return false
}

func NICRates(interval time.Duration) ([]NicRate, error) {
	// Snapshot 1
	ifStats1, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	start := time.Now()
	time.Sleep(interval)
	elapsed := time.Since(start).Seconds()

	// Snapshot 2
	ifStats2, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	// Indexer snapshot 1 par nom d’interface
	s1 := make(map[string]net.IOCountersStat, len(ifStats1))
	for _, s := range ifStats1 {
		s1[s.Name] = s
	}

	// Indexer meta (MTU, adresses, flags) par nom
	meta := make(map[string]net.InterfaceStat, len(ifaces))
	for _, inf := range ifaces {
		meta[inf.Name] = inf
	}

	out := make([]NicRate, 0, len(ifStats2))
	for _, s2 := range ifStats2 {
		inf := meta[s2.Name]
		// Filtrer : UP et non-loopback
		if !isUp(inf.Flags) || isLoopback(inf.Flags) {
			continue
		}

		s1prev, ok := s1[s2.Name]
		if !ok {
			continue // interface apparue entre temps
		}
		rxBps := float64(s2.BytesRecv-s1prev.BytesRecv) / elapsed
		txBps := float64(s2.BytesSent-s1prev.BytesSent) / elapsed

		n := NicRate{
			Name:   s2.Name,
			Up:     true,
			MTU:    int(inf.MTU),
			RxBps:  rxBps,
			TxBps:  txBps,
			RxMbps: rxBps * 8 / 1_000_000, // Mb/s (decimal)
			TxMbps: txBps * 8 / 1_000_000,
		}
		for _, a := range inf.Addrs {
			n.Addr = append(n.Addr, a.Addr)
		}
		out = append(out, n)
	}
	return out, nil
}
