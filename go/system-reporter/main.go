package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	cpu "github.com/shirou/gopsutil/v3/cpu"
	disk "github.com/shirou/gopsutil/v3/disk"
	host "github.com/shirou/gopsutil/v3/host"
	mem "github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func main() {
	fmt.Println("--- System Information Report ---")

	// Host and OS Info
	hostInfo, err := host.Info()
	if err != nil {
		log.Printf("Error getting host info: %v", err)
	} else {
		fmt.Println("\n--- Host & OS ---")
		fmt.Printf("Hostname: %s\n", hostInfo.Hostname)
		fmt.Printf("OS: %s\n", hostInfo.OS)
		fmt.Printf("Platform: %s\n", hostInfo.Platform)
		fmt.Printf("Platform Family: %s\n", hostInfo.PlatformFamily)
		fmt.Printf("Platform Version: %s\n", hostInfo.PlatformVersion)
		fmt.Printf("Kernel Version: %s\n", hostInfo.KernelVersion)
		fmt.Printf("Architecture: %s\n", hostInfo.KernelArch)
		fmt.Printf("Virtualization System: %s\n", hostInfo.VirtualizationSystem)
		fmt.Printf("Virtualization Role: %s\n", hostInfo.VirtualizationRole)
		fmt.Printf("Uptime: %v\n", formatDuration(time.Duration(hostInfo.Uptime)*time.Second))
		fmt.Printf("Boot Time: %s\n", time.Unix(int64(hostInfo.BootTime), 0).Format("2006-01-02 15:04:05"))
	}

	fmt.Println("---------------------------------")

	// CPU Info
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Printf("Error getting CPU info: %v", err)
	} else {
		fmt.Println("\n--- CPU ---")
        logicalCores, err := cpu.Counts(true)
        if err != nil {
            log.Printf("Error getting logical CPU count: %v", err)
        } else {
            fmt.Printf("  Total Logical Cores: %d\n", logicalCores)
        }
        for _, cpu := range cpuInfo {
            fmt.Printf("  Model Name: %s\n", cpu.ModelName)
            fmt.Printf("  Physical Cores: %d\n", cpu.Cores)
            fmt.Printf("  Mhz: %.2f\n", cpu.Mhz)
            fmt.Printf("  Cache Size: %d KB\n", cpu.CacheSize)
        }
	}

	// Memory Info
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting virtual memory info: %v", err)
	} else {
		fmt.Println("\n--- Memory ---")
		fmt.Printf("  Total: %s\n", formatBytes(vmStat.Total))
		fmt.Printf("  Used: %s\n", formatBytes(vmStat.Used))
		fmt.Printf("  Free: %s\n", formatBytes(vmStat.Free))
		fmt.Printf("  Used Percent: %.2f%%\n", vmStat.UsedPercent)
	}

	swapStat, err := mem.SwapMemory()
	if err != nil {
		log.Printf("Error getting swap memory info: %v", err)
	} else {
		fmt.Printf("  Swap Total: %s\n", formatBytes(swapStat.Total))
		fmt.Printf("  Swap Used: %s\n", formatBytes(swapStat.Used))
		fmt.Printf("  Swap Free: %s\n", formatBytes(swapStat.Free))
		fmt.Printf("  Swap Used Percent: %.2f%%\n", swapStat.UsedPercent)
	}

	// Disk Info
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Error getting disk partitions: %v", err)
	} else {
		fmt.Println("\n--- Disk Usage ---")
		for _, p := range partitions {
			usage, err := disk.Usage(p.Mountpoint)
			if err != nil {
				log.Printf("Error getting disk usage for %s: %v", p.Mountpoint, err)
				continue
			}
			fmt.Printf("  %s (%s): Total %s, Used %s, Free %s (%.2f%% used)\n",
				p.Mountpoint, p.Fstype, formatBytes(usage.Total), formatBytes(usage.Used),
				formatBytes(usage.Free), usage.UsedPercent)
		}
	}

	// Network Info
	netInfo, err := net.Interfaces()
	if err != nil {
		log.Printf("Error getting network interfaces: %v", err)
	} else {
		fmt.Println("\n--- Network Interfaces ---")
		for _, ni := range netInfo {
			fmt.Printf("  Name: %s\n", ni.Name)
			fmt.Printf("    MAC Address: %s\n", ni.HardwareAddr)
			fmt.Printf("    Flags: %v\n", ni.Flags)
			for _, addr := range ni.Addrs {
				fmt.Printf("    IP Address: %s\n", addr.Addr)
			}
			// Add more network stats if needed, e.g., bytes sent/received
			// ioCounters, err := net.IOCounters(true)
			// if err == nil {
			// 	for _, io := range ioCounters {
			// 		if io.Name == ni.Name {
			// 			fmt.Printf("    Bytes Sent: %s, Bytes Recv: %s\n", formatBytes(io.BytesSent), formatBytes(io.BytesRecv))
			// 		}
			// 	}
			// }
		}
	}
}

// formatBytes formats bytes into a human-readable string (e.g., KB, MB, GB).
func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// formatDuration formats a duration into a human-readable string.
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	parts := []string{}
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%dh", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%dm", m))
	}
	if s > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", s))
	}
	return strings.Join(parts, " ")
}
