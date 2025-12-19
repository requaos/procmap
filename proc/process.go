package proc

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo holds information about a single process
type ProcessInfo struct {
	PID         int32
	Name        string
	CPUPercent  float64
	MemPercent  float32
	Status      string
	CreateTime  time.Time
	NumThreads  int32
	Username    string
}

// GetProcesses returns a list of all running processes with their details
func GetProcesses() ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var processInfos []ProcessInfo
	for _, p := range procs {
		info, err := getProcessInfo(p)
		if err != nil {
			// Skip processes we can't read (permission denied, etc.)
			continue
		}
		processInfos = append(processInfos, info)
	}

	return processInfos, nil
}

func getProcessInfo(p *process.Process) (ProcessInfo, error) {
	name, _ := p.Name()
	cpuPercent, _ := p.CPUPercent()
	memPercent, _ := p.MemoryPercent()
	status, _ := p.Status()
	createTime, _ := p.CreateTime()
	numThreads, _ := p.NumThreads()
	username, _ := p.Username()

	return ProcessInfo{
		PID:        p.Pid,
		Name:       name,
		CPUPercent: cpuPercent,
		MemPercent: memPercent,
		Status:     status[0],
		CreateTime: time.Unix(createTime/1000, 0),
		NumThreads: numThreads,
		Username:   username,
	}, nil
}
