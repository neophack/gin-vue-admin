package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// These variables are set in build step
var (
	Version  = "unset"
	Revision = "unset"
)

const (
	memFreeRatio     = 5
	gpuFreeRatio     = 5
	memModerateRatio = 90
	gpuModerateRatio = 75
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type Server struct {
	Os   Os   `json:"os"`
	Cpu  Cpu  `json:"cpu"`
	Ram  Ram  `json:"ram"`
	Disk Disk `json:"disk"`
	Gpu  Gpu  `json:"gpu"`
}

type Os struct {
	GOOS         string `json:"goos"`
	NumCPU       int    `json:"numCpu"`
	Compiler     string `json:"compiler"`
	GoVersion    string `json:"goVersion"`
	NumGoroutine int    `json:"numGoroutine"`
}

type Cpu struct {
	Cpus  []float64 `json:"cpus"`
	Cores int       `json:"cores"`
}

type Ram struct {
	UsedMB      int `json:"usedMb"`
	TotalMB     int `json:"totalMb"`
	UsedPercent int `json:"usedPercent"`
}

type Disk struct {
	UsedMB      int `json:"usedMb"`
	UsedGB      int `json:"usedGb"`
	TotalMB     int `json:"totalMb"`
	TotalGB     int `json:"totalGb"`
	UsedPercent int `json:"usedPercent"`
}

type GpuInfo struct {
	Index           int    `json:"index"`
	GpuUUID         string `json:"gpu_uuid"`
	Name            string `json:"name"`
	MemoryUsed      int    `json:"memory_used"`
	MemoryTotal     int    `json:"memory_total"`
	UtilizationGpu  int    `json:"utilization_gpu"`
	PersistanceMode bool   `json:"persistance_mode"`
}

type Process struct {
	GpuUUID       string `json:"gpu_uuid"`
	Pid           int    `json:"pid"`
	UsedGpuMemory int    `json:"used_gpu_memory"`
	User          string `json:"user"`
	Command       string `json:"command"`
}
type Gpu struct {
	GpuInfos  []GpuInfo `json:"gpu_infos"`
	Processes []Process `json:"processes"`
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: InitCPU
//@description: OS信息
//@return: o Os, err error

func InitOS() (o Os) {
	o.GOOS = runtime.GOOS
	o.NumCPU = runtime.NumCPU()
	o.Compiler = runtime.Compiler
	o.GoVersion = runtime.Version()
	o.NumGoroutine = runtime.NumGoroutine()
	return o
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: InitCPU
//@description: CPU信息
//@return: c Cpu, err error

func InitCPU() (c Cpu, err error) {
	if cores, err := cpu.Counts(false); err != nil {
		return c, err
	} else {
		c.Cores = cores
	}
	if cpus, err := cpu.Percent(time.Duration(200)*time.Millisecond, true); err != nil {
		return c, err
	} else {
		c.Cpus = cpus
	}
	return c, nil
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: InitRAM
//@description: RAM信息
//@return: r Ram, err error

func InitRAM() (r Ram, err error) {
	if u, err := mem.VirtualMemory(); err != nil {
		return r, err
	} else {
		r.UsedMB = int(u.Used) / MB
		r.TotalMB = int(u.Total) / MB
		r.UsedPercent = int(u.UsedPercent)
	}
	return r, nil
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: InitDisk
//@description: 硬盘信息
//@return: d Disk, err error

func InitDisk() (d Disk, err error) {
	if u, err := disk.Usage("/"); err != nil {
		return d, err
	} else {
		d.UsedMB = int(u.Used) / MB
		d.UsedGB = int(u.Used) / GB
		d.TotalMB = int(u.Total) / MB
		d.TotalGB = int(u.Total) / GB
		d.UsedPercent = int(u.UsedPercent)
	}
	return d, nil
}

func GetUserFromPid(pid int) string {
	out, err := exec.Command("ps", "ho", "user", strconv.Itoa(pid)).Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(out), "\n")
}

func GetCommandFromPid(pid int) string {
	out, err := exec.Command("ps", "ho", "command", strconv.Itoa(pid)).Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(out), "\n")
}

func NewGpuInfoFromLine(line string) GpuInfo {
	values := strings.Split(line, ", ")

	index, err := strconv.Atoi(values[0])
	if err != nil {
		log.Fatal(err)
	}
	gpuUUID := values[1]
	name := values[2]
	memoryUsed, err := strconv.Atoi(values[3])
	if err != nil {
		log.Fatal(err)
	}
	memoryTotal, err := strconv.Atoi(values[4])
	if err != nil {
		log.Fatal(err)
	}
	utilizationGpu, err := strconv.Atoi(values[5])
	if err != nil {
		log.Fatal(err)
	}
	persistanceMode := values[6]
	return GpuInfo{
		Index:           index,
		GpuUUID:         gpuUUID,
		Name:            name,
		MemoryUsed:      memoryUsed,
		MemoryTotal:     memoryTotal,
		UtilizationGpu:  utilizationGpu,
		PersistanceMode: persistanceMode == "Enabled",
	}
}

func NewProcessFromLine(line string) Process {
	values := strings.Split(line, ", ")

	gpuUUID := values[0]
	pid, err := strconv.Atoi(values[1])
	if err != nil {
		log.Fatal(err)
	}
	user := GetUserFromPid(pid)
	usedGpuMemory, err := strconv.Atoi(values[2])
	if err != nil {
		log.Fatal(err)
	}
	command := GetCommandFromPid(pid)

	return Process{GpuUUID: gpuUUID, Pid: pid, UsedGpuMemory: usedGpuMemory, User: user, Command: command}

}

func RetrieveGpus() map[string]GpuInfo {
	out, err := exec.Command(
		"/usr/bin/env", "nvidia-smi",
		"--format=csv,noheader,nounits",
		"--query-gpu=index,gpu_uuid,name,memory.used,memory.total,utilization.gpu,persistence_mode").Output()

	if err != nil {
		//log.Fatal(err)
		fmt.Println("Cant find any gpu!")
		return make(map[string]GpuInfo, 10)
	}

	outStr := strings.TrimSuffix(string(out), "\n")
	lines := strings.Split(outStr, "\n")

	gpus := make(map[string]GpuInfo, 10)
	for _, line := range lines {
		gpu := NewGpuInfoFromLine(line)
		gpus[gpu.GpuUUID] = gpu
	}
	return gpus
}

func RetrieveProcesses() []Process {
	out, err := exec.Command(
		"/usr/bin/env", "nvidia-smi",
		"--format=csv,noheader,nounits",
		"--query-compute-apps=gpu_uuid,pid,used_memory",
	).Output()
	if err != nil {
		//log.Fatal(err)
		fmt.Println("Cant find any process!")
		return []Process{}
	}

	outStr := strings.TrimSuffix(string(out), "\n")
	lines := strings.Split(outStr, "\n")
	if lines[0] == "" {
		return []Process{}
	}

	processes := []Process{}
	for _, line := range lines {
		process := NewProcessFromLine(line)
		processes = append(processes, process)
	}
	return processes
}

func gpuProcessExists(gpu GpuInfo, processes []Process) string {
	for _, process := range processes {
		if gpu.GpuUUID == process.GpuUUID {
			return "RUNNING"
		}
	}
	return "-------"
}

func printProcesses(processes []Process, gpus map[string]GpuInfo) string {
	outputs := []string{}
	for _, process := range processes {
		outputs = append(outputs, fmt.Sprintf("| %3d | %10s | %7d | %5d MiB | %22.22s |",
			gpus[process.GpuUUID].Index,
			process.User,
			process.Pid,
			process.UsedGpuMemory,
			process.Command))
	}
	return strings.Join(outputs, "\n")
}

func printWarnPersistanceMode(gpus map[string]GpuInfo) {
	for k := range gpus {
		gpu := gpus[k]
		if !gpu.PersistanceMode {
			fmt.Println("Consider enabling persistence mode on your Gpu(s) for faster response.")
			fmt.Println("For more information: https://docs.nvidia.com/deploy/driver-persistence/")
			break
		}
	}
}

func sortByGpuInfoIndex(msg map[string]GpuInfo) []GpuInfo {
	mis := map[int]string{}
	miskeys := []int{}
	for k, v := range msg {
		mis[v.Index] = k
		miskeys = append(miskeys, v.Index)
	}
	sort.Ints(miskeys)

	gpus := make([]GpuInfo, 0, len(miskeys))
	for _, v := range miskeys {
		gpuUUID := mis[v]
		gpus = append(gpus, msg[gpuUUID])
	}
	return gpus
}

func printWithColor(gpu GpuInfo, processes []Process) {
	usedMem := gpu.MemoryUsed
	totalMem := gpu.MemoryTotal
	gpuUtil := gpu.UtilizationGpu
	memUtil := usedMem / totalMem
	isModerate := false
	isHigh := float32(gpuUtil) >= gpuModerateRatio || float32(memUtil) >= memModerateRatio
	if !isHigh {
		isModerate = float32(gpuUtil) >= gpuFreeRatio || float32(memUtil) >= memFreeRatio
	}

	var au func(interface{}) aurora.Value
	if isHigh {
		au = aurora.Red
	} else if isModerate {
		au = aurora.Yellow
	} else {
		au = aurora.Green
	}

	colorFormat := "| %3d %22s | %3d  | %s | %s |\n"
	fmt.Printf(
		colorFormat,
		au(gpu.Index),
		au(gpu.Name),
		au(gpu.UtilizationGpu),
		au(fmt.Sprintf("%5d / %5d MiB", gpu.MemoryUsed, gpu.MemoryTotal)),
		au(gpuProcessExists(gpu, processes)),
	)
}

func InitGpu() (d Gpu, err error) {
	gpus := RetrieveGpus()
	processes := RetrieveProcesses()
	sortedGpus := sortByGpuInfoIndex(gpus)
	d.GpuInfos = sortedGpus
	d.Processes = processes

	return d, nil
}
