package main

import (
	"errors"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcDTO struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name,omitempty"`
	Cmdline       string  `json:"cmdline,omitempty"`
	Status        string  `json:"status,omitempty"`
	CreateTime    int64   `json:"create_time,omitempty"` // ms epoch
	Username      string  `json:"username,omitempty"`
	NumThreads    int32   `json:"threads,omitempty"`
	MemoryRSS     uint64  `json:"memory_rss,omitempty"`
	MemoryVMS     uint64  `json:"memory_vms,omitempty"`
	MemoryPercent float32 `json:"memory_percent,omitempty"`
}

// Data Transfert Object pour Procs
func DTOProc(p *process.Process) (ProcDTO, error) {
	dto := ProcDTO{PID: p.Pid}

	if name, err := p.Name(); err == nil {
		dto.Name = name
	}
	if cl, err := p.Cmdline(); err == nil {
		dto.Cmdline = cl
	}

	if st, err := p.Status(); err == nil {
		stat := strings.ToUpper(st[0])
		dto.Status = strings.ToUpper(stat[:1])
	}

	if ct, err := p.CreateTime(); err == nil {
		dto.CreateTime = ct
	}
	if u, err := p.Username(); err == nil {
		dto.Username = u
	}
	if th, err := p.NumThreads(); err == nil {
		dto.NumThreads = th
	}
	if pm, err := p.MemoryInfo(); err == nil {
		dto.MemoryRSS = pm.RSS
		dto.MemoryVMS = pm.VMS
	}
	if mp, err := p.MemoryPercent(); err == nil {
		dto.MemoryPercent = mp
	}
	return dto, nil
}

// Récupération de tous les processus
func DTOProcAllLoad() (*[]ProcDTO, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	out := make([]ProcDTO, 0, len(procs))
	for _, p := range procs {
		dto, _ := DTOProc(p)
		out = append(out, dto)
	}
	return &out, nil
}

// Récupération du processus id à partir Procs
func DTOProcLoad(id int32) (*ProcDTO, error) {
	procs := *Datas.Procs
	for _, p := range procs {
		if p.PID == id {
			return &p, nil
		}
	}
	err := errors.New("processus inconnu")
	return nil, err
}
