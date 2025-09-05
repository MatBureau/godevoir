package moninfluxdb

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/shirou/gopsutil/v4/load"
)

func Open(host, name, token string) (*influxdb3.Client, error) {
	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     host,
		Token:    token,
		Database: name,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func WriteLoad(client *influxdb3.Client, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var mstats load.AvgStat
	if err := json.NewDecoder(resp.Body).Decode(&mstats); err != nil {
		return err
	}

	pt := influxdb3.NewPoint(
		"load_avg",                             // measurement
		map[string]string{"host": "localhost"}, // tags
		map[string]any{ // fields
			"load1":  mstats.Load1,
			"load5":  mstats.Load5,
			"load15": mstats.Load15,
		},
		time.Now().UTC(), // timestamp d’échantillonnage
	)
	if err := client.WritePoints(context.Background(), []*influxdb3.Point{pt}); err != nil {
		return err
	}
	return nil
}

type MemDTO struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Free        uint64  `json:"free"`
}

func WriteMemFromURL(client *influxdb3.Client, url, host string) error {
	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var m MemDTO
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return err
	}

	pt := influxdb3.NewPoint(
		"mem",
		map[string]string{"host": host},
		map[string]any{
			"total":        int64(m.Total),
			"available":    int64(m.Available),
			"used":         int64(m.Used),
			"free":         int64(m.Free),
			"used_percent": m.UsedPercent,
		},
		time.Now().UTC(),
	)
	return client.WritePoints(ctx, []*influxdb3.Point{pt})
}

// + imports: "strconv"
type cpuPercentsDTO []float64

func WriteCPUPercent(client *influxdb3.Client, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var percs cpuPercentsDTO
	if err := json.NewDecoder(resp.Body).Decode(&percs); err != nil {
		return err
	}

	pts := make([]*influxdb3.Point, 0, len(percs))
	for i, v := range percs {
		pt := influxdb3.NewPoint(
			"cpu_percent",
			map[string]string{"host": "localhost", "cpu": strconv.Itoa(i)},
			map[string]any{"value": v},
			time.Now().UTC(),
		)
		pts = append(pts, pt)
	}
	return client.WritePoints(context.Background(), pts)
}

type NicRateDTO struct {
	Name   string  `json:"name"`
	RxBps  float64 `json:"rx_bps"`
	TxBps  float64 `json:"tx_bps"`
	RxMbps float64 `json:"rx_mbps"`
	TxMbps float64 `json:"tx_mbps"`
	MTU    int     `json:"mtu"`
	Up     bool    `json:"up"`
}

func WriteNics(client *influxdb3.Client, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var nics []NicRateDTO
	if err := json.NewDecoder(resp.Body).Decode(&nics); err != nil {
		return err
	}

	pts := make([]*influxdb3.Point, 0, len(nics))
	for _, n := range nics {
		pt := influxdb3.NewPoint(
			"nic_rate",
			map[string]string{"host": "localhost", "nic": n.Name},
			map[string]any{
				"rx_bps":  n.RxBps,
				"tx_bps":  n.TxBps,
				"rx_mbps": n.RxMbps,
				"tx_mbps": n.TxMbps,
				"mtu":     int64(n.MTU),
				"up":      n.Up,
			},
			time.Now().UTC(),
		)
		pts = append(pts, pt)
	}
	return client.WritePoints(context.Background(), pts)
}

type ProcDTO struct {
	Status string `json:"status"`
}

func WriteProcsCount(client *influxdb3.Client, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var procs []ProcDTO
	if err := json.NewDecoder(resp.Body).Decode(&procs); err != nil {
		return err
	}

	counts := map[string]int64{}
	for _, p := range procs {
		s := p.Status
		if s == "" {
			s = "UNK"
		}
		counts[s]++
	}
	// total
	var total int64
	for _, c := range counts {
		total += c
	}

	pts := make([]*influxdb3.Point, 0, len(counts)+1)
	for s, c := range counts {
		pts = append(pts, influxdb3.NewPoint(
			"procs_count",
			map[string]string{"host": "localhost", "status": s},
			map[string]any{"count": c},
			time.Now().UTC(),
		))
	}
	pts = append(pts, influxdb3.NewPoint(
		"procs_count",
		map[string]string{"host": "localhost", "status": "TOTAL"},
		map[string]any{"count": total},
		time.Now().UTC(),
	))
	return client.WritePoints(context.Background(), pts)
}

type DiskUsageDTO struct {
	Path        string  `json:"path"`
	FSType      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

func WriteDiskUsage(client *influxdb3.Client, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var items []DiskUsageDTO
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return err
	}

	pts := make([]*influxdb3.Point, 0, len(items))
	now := time.Now().UTC()
	for _, d := range items {
		pts = append(pts, influxdb3.NewPoint(
			"disk_usage",
			map[string]string{"host": "localhost", "mountpoint": d.Path, "fstype": d.FSType},
			map[string]any{
				"total":        int64(d.Total),
				"free":         int64(d.Free),
				"used":         int64(d.Used),
				"used_percent": d.UsedPercent,
			},
			now,
		))
	}
	return client.WritePoints(context.Background(), pts)
}
