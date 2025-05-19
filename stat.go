package main

import (
	"fmt"
	"time"

	"sort"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {

	for {

		cpuPercent, _ := cpu.Percent(0, false)

		// Get memory usage
		vmStat, _ := mem.VirtualMemory()

		diskStat, _ := disk.Usage("/")

		proc, _ := process.Processes()

		fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent[0])
		fmt.Printf("Memory Usage: %.2fGB Free %.2f%% of %.2f GB\n", float64(vmStat.Free)/(1024*1024*1024), vmStat.UsedPercent, float64(vmStat.Total)/(1024*1024*1024))
		fmt.Printf("Disk Used: %.2f%%  %.2fGB  Free: %.2fGB Total: %.2f\n", diskStat.UsedPercent, float64(diskStat.Used)/(1024*1024*1024), float64(diskStat.Free)/(1024*1024*1024), float64(diskStat.Total)/(1024*1024*1024))

		type procCPU struct {
			name string
			cpu  float64
		}

		type procMem struct {
			name string
			mem  float32
		}

		var procList []procCPU
		var memList []procMem

		for _, p := range proc {
			vmStat, err := p.MemoryPercent()
			if err != nil {
				continue
			}
			name, err := p.Name()
			if err != nil {
				continue
			}
			memList = append(memList, procMem{name: name, mem: vmStat})
		}

		sort.Slice(memList, func(i, j int) bool {
			return memList[i].mem > memList[j].mem
		})
		fmt.Println("\n Memory Usage:\n")
		for i := 0; i < 5 && i < len(memList); i++ {

			fmt.Printf("\n %d. %s: %.2f%% Memory\n", i+1, memList[i].name, memList[i].mem)
		}

		for _, p := range proc {
			cpuPercent, err := p.CPUPercent()
			if err != nil {
				continue
			}
			name, err := p.Name()
			if err != nil {
				continue
			}
			procList = append(procList, procCPU{name: name, cpu: cpuPercent})
		}

		sort.Slice(procList, func(i, j int) bool {
			return procList[i].cpu > procList[j].cpu
		})
		fmt.Println("\nCPU Usage:\n")

		for i := 0; i < 5 && i < len(procList); i++ {

			fmt.Printf("\n %d. %s: %.2f%% CPU\n", i+1, procList[i].name, procList[i].cpu)
		}

		// Wait 2 seconds before checking again
		fmt.Println("------")
		time.Sleep(2 * time.Second)
	}

}
