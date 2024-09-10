package info

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/distatus/battery"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func GetCPUTemperature() string {
	cmd := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to get cpu temp :%v", err)
	}
	tempStr := strings.TrimSpace(string(output))

	// convert the temperature string to an integer
	tempInt, err := strconv.Atoi(tempStr)
	if err != nil {
		log.Fatalf("Failed to convert temp to integer :%v", err)
	}
	tempCelcius := float64(tempInt) / 1000.0
	// Print the temperature with the degree symbol and "C"
	return fmt.Sprintf(" %.1f°C\n", tempCelcius)
}

func GetDiskInfo() string {
	usage, err := disk.Usage("/")
	if err != nil {
		return "N/A"
	}

	total := float64(usage.Total) / (1024 * 1024 * 1024) // Convert to GiB
	used := float64(usage.Used) / (1024 * 1024 * 1024)   // Convert to GiB

	return fmt.Sprintf("%.2f GiB used / %.2f GiB total (%.0f%% used)", used, total, usage.UsedPercent)
}

func GetOSInfo() string {
	hostInfo, _ := host.Info()
	return fmt.Sprintf("%s %s %s", hostInfo.Platform, hostInfo.PlatformVersion, runtime.GOARCH)
}

func GetMachineInfo() string {
	hostInfo, _ := host.Info()
	return hostInfo.Hostname
}

func GetKernelInfo() string {
	hostInfo, _ := host.Info()
	return hostInfo.KernelVersion
}

func GetUptime() string {
	hostInfo, _ := host.Info()
	uptime := time.Duration(hostInfo.Uptime) * time.Second
	return uptime.Round(time.Minute).String()
}

func GetShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "Unknown"
	}
	return filepath.Base(shell)
}

func GetCPUInfo() string {
	cpuInfo, _ := cpu.Info()
	if len(cpuInfo) > 0 {
		return fmt.Sprintf("%s (%d) @ %.2f GHz", cpuInfo[0].ModelName, runtime.NumCPU(), cpuInfo[0].Mhz/1000)
	}
	return "Unknown"
}

func GetMemoryInfo() string {
	v, _ := mem.VirtualMemory()
	return fmt.Sprintf("%.2f GiB / %.2f GiB (%.0f%%)", float64(v.Used)/(1024*1024*1024), float64(v.Total)/(1024*1024*1024), v.UsedPercent)
}

func GetBatteryInfo() string {
	batteries, err := battery.GetAll()
	if err != nil || len(batteries) == 0 {
		return "N/A"
	}

	// We'll just use the first battery
	bat := batteries[0]
	percentage := bat.Current / bat.Full * 100

	return fmt.Sprintf("%.0f%%", percentage)
}
