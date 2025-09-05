package moninfluxdb

import (
	"context"
	"encoding/json"
	"net/http"
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
