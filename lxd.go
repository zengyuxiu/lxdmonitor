package main

import (
	"fmt"
	lxd "github.com/lxc/lxd/client"
	"github.com/struCoder/pidusage"
	"math"
	"time"
)

type Netinfo struct {
	BytesReceived   int64
	BytesSent       int64
	PacketsReceived int64
	PacketsSent     int64
}

type Sourceinfo struct {
	CpuUsage float64
	MemUsage float64
}

func Getnetstat(instance string, nic string, client lxd.InstanceServer) Netinfo {
	state, _, err := client.GetInstanceState(instance)
	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}
	return Netinfo(state.Network[nic].Counters)
}
func Getsource(instance string, client lxd.InstanceServer) Sourceinfo {
	var cpuusage float64
	var memusage float64
	current_instance, _, err := client.GetInstance(instance)
	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}
	if current_instance.Type == "container" {
		state, _, _ := client.GetInstanceState(instance)
		memusage = float64(state.Memory.Usage)
		cpu := state.Memory.Usage
		checktime := time.Now()
		state, _, _ = client.GetInstanceState(instance)
		cpu_current := state.CPU.Usage
		dur := time.Now().Sub(checktime)
		elapsed_cpu := cpu_current - cpu
		cpuusage = float64(elapsed_cpu) / float64(dur) * math.Pow10(-3) // per ns
	} else if current_instance.Type == "virtual-machine" {
		state, _, _ := client.GetInstanceState(instance)
		sysInfo, err := pidusage.GetStat(int(state.Pid))
		if err != nil {
			fmt.Printf(err.Error())
			panic(err)
		}
		memusage = sysInfo.Memory
		cpuusage = sysInfo.CPU
		/*		cpu := int64(sysInfo.CPU)
				checktime := time.Now()
				sysInfo, err = pidusage.GetStat(int(state.Pid))
				if err != nil{
					fmt.Printf(err.Error())
					panic(err)
				}
				cpu_current := int64(sysInfo.CPU)
				dur := time.Now().Sub(checktime)
				elapsed_cpu := cpu_current-cpu
				cpuusage = float64(elapsed_cpu * 100 / int64(dur))*/
	}
	return Sourceinfo{
		CpuUsage: cpuusage,
		MemUsage: memusage,
	}
}
