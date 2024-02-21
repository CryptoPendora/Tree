//please donate 12CsSsSv7UZuoxZJPzyfrSrwXfK1SbpkVM

package main

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

func displayHostInfo() {
	info, _ := host.Info()
	fmt.Printf("Hostname: %s\n", info.Hostname)
	fmt.Printf("OS: %s\n", info.OS)
	fmt.Printf("Platform: %s\n", info.Platform)
	fmt.Printf("Platform Version: %s\n", info.PlatformVersion)
	fmt.Printf("Uptime: %d sec\n", info.Uptime)
}

func displayLoadAverage() {
	avg, _ := load.Avg()
	fmt.Printf("Load Average 1 min: %.2f\n", avg.Load1)
	fmt.Printf("Load Average 5 min: %.2f\n", avg.Load5)
	fmt.Printf("Load Average 15 min: %.2f\n", avg.Load15)
}

func displayNetworkStats() {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		fmt.Printf("Interface Name: %s\n", inter.Name)
		fmt.Printf("Hardware Address: %s\n", inter.HardwareAddr)
		fmt.Printf("Flags: %v\n", inter.Flags)
	}
}

func displayNetworkIOStats() {
	stats, err := net.IOCounters(true)
	if err != nil {
		fmt.Println("Error retrieving network IO stats:", err)
		return
	}
	for _, stat := range stats {
		fmt.Printf("Interface: %v, Bytes Sent: %v, Bytes Received: %v\n", stat.Name, stat.BytesSent, stat.BytesRecv)
	}
}

func displayProcessInfo() {
	processes, _ := process.Processes()
	for _, p := range processes[:5] { // Limit to first 5 processes to keep output short
		name, _ := p.Name()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		fmt.Printf("Process Name: %s, CPU: %.2f%%, Memory: %.2f%%\n", name, cpuPercent, memPercent)
	}
}

func displayDiskUsage() {
	usages, err := disk.Partitions(true)
	if err != nil {
		fmt.Printf("Error retrieving disk partitions: %v\n", err)
		return
	}
	for _, partition := range usages {
		usageStat, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			fmt.Printf("Error retrieving disk usage: %v\n", err)
			continue
		}
		fmt.Printf("Disk: %s, Mount: %s, Total: %.2fGB, Used: %.2fGB, Free: %.2fGB, Usage: %.2f%%\n",
			partition.Device, partition.Mountpoint,
			float64(usageStat.Total)/float64(1<<30),
			float64(usageStat.Used)/float64(1<<30),
			float64(usageStat.Free)/float64(1<<30),
			usageStat.UsedPercent)
	}
}

func displayCPUInfo() {
	infos, err := cpu.Info()
	if err != nil {
		fmt.Printf("Error retrieving CPU info: %v\n", err)
		return
	}
	for _, info := range infos {
		fmt.Printf("CPU: %v, Cores: %v, Model Name: %s, Speed: %.2fMHz\n", info.PhysicalID, info.Cores, info.ModelName, info.Mhz)
	}
	// แสดงเปอร์เซ็นต์การใช้งาน
	percent, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Printf("Error retrieving CPU percent: %v\n", err)
		return
	}
	fmt.Printf("CPU Usage: %.2f%%\n", percent[0])
}

func displayMemoryInfo() {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("Error retrieving VM info: %v\n", err)
		return
	}
	fmt.Printf("Total RAM: %.2fGB\n", float64(vmStat.Total)/float64(1<<30))
	fmt.Printf("Used RAM: %.2fGB (%.2f%%)\n", float64(vmStat.Used)/float64(1<<30), vmStat.UsedPercent)
	fmt.Printf("Free RAM: %.2fGB\n", float64(vmStat.Available)/float64(1<<30))
}

func main() {
	displayHostInfo()
	displayLoadAverage()
	displayNetworkStats()
	displayNetworkIOStats()
	displayProcessInfo()
	displayCPUInfo()
	displayDiskUsage()
	displayMemoryInfo()

}
