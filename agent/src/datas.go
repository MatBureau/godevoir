package main

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

type AllDatas struct {
	Load    *load.AvgStat
	CPU     *[]cpu.InfoStat
	CPULoad *[]float64
	Parts   *[]disk.PartitionStat
	Procs   *[]ProcDTO
	Nics    *[]NicRate
	Mem     *mem.VirtualMemoryStat
}

var Datas AllDatas
