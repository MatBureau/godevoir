package main

import (
	"client/moninfluxdb"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskUsageDTO struct {
	Path        string  `json:"path"`
	FSType      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

func goDiskUsage(client *influxdb3.Client) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		parts, _ := disk.Partitions(false)
		items := make([]DiskUsageDTO, 0, len(parts))
		for _, p := range parts {
			if u, err := disk.Usage(p.Mountpoint); err == nil {
				items = append(items, DiskUsageDTO{
					Path:        p.Mountpoint,
					FSType:      p.Fstype,
					Total:       u.Total,
					Free:        u.Free,
					Used:        u.Used,
					UsedPercent: u.UsedPercent,
				})
			}
		}
		Datas.DiskUsage = &items
		_ = moninfluxdb.WriteDiskUsage(client, ServerURL+"/disks/usage")
		LogMessage("goroutine: goDiskUsage")
	}
}
