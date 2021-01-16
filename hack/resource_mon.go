package main

import (
        "fmt"
        "github.com/shirou/gopsutil/mem"
        "github.com/shirou/gopsutil/process"
        "strconv"
)

func main() {
        pids, _ := process.Pids()
        for _, pid := range pids {
                fmt.Println("pid: ", pid)
                proc, _ := process.NewProcess(pid)
                username, _ := proc.Username()
                fmt.Println("username: ", username)
                ppid, _ := proc.Ppid()
                fmt.Println("ppid: ", ppid)
                name, _ := proc.Name()
                fmt.Println("name: ", name)
                cmdline, _ := proc.Cmdline()
                fmt.Println("cmdline: ", cmdline)
                meminfo, _ := proc.MemoryInfo()
                fmt.Println("rss: ", meminfo.RSS)
                fmt.Println("vms: ", meminfo.VMS)
                fmt.Println("swap: ", meminfo.Swap)
                memory_maps, err := proc.MemoryMaps(true)
                if err != nil {
                        fmt.Println("Error: ", err)
                        continue
                }
                pss := uint64(0)
                uss := uint64(0)
                for _, memory_map := range *memory_maps {
                        pss += memory_map.Pss
                        uss += memory_map.PrivateClean + memory_map.PrivateDirty
                }
                fmt.Println("pss: ", pss)
                fmt.Println("uss: ", uss)
                num_fds, _ := proc.NumFDs()
                fmt.Println("num fds: ", num_fds)
                fmt.Println("-----")
        }
        vmStat, _ := mem.VirtualMemory()
        fmt.Println("Total memory: " + strconv.FormatUint(vmStat.Total, 10) + "bytes")
        fmt.Println("Free memory: " + strconv.FormatUint(vmStat.Free, 10) + "bytes")
        fmt.Println("Percentage used memory: " + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%")
}
