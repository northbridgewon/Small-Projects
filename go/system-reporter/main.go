package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	host "github.com/shirou/gopsutil/v3/host"
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
