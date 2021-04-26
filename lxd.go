package main

import (
	"fmt"
	"github.com/struCoder/pidusage"
	"time"
)
type Netinfo struct {
	BytesReceived   int64
	BytesSent       int64
	PacketsReceived int64
	PacketsSent     int64
}

type Sourceinfo struct {
	CpuUsage        int64
	MemUsage        int64
}

func Getnetstat(instance string,nic string) Netinfo {
	d, err :=InitLxdInstanceServer()
	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}
	c := *d
	state,_,err := c.GetInstanceState(instance)
	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}
	return Netinfo(state.Network[nic].Counters)
}
func Getsource(instance string,instancetype string) Sourceinfo {
	var cpuusage int64
	var memusage int64
	d, err :=InitLxdInstanceServer()
	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}
	c := *d
	_,_,err = c.GetInstance(instance)
	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}
	if instancetype == "container"{
		state,_,_ := c.GetInstanceState(instance)
		memusage = state.Memory.Usage
		cpu := state.CPU.Usage
		checktime := time.Now()
		state,_,_ = c.GetInstanceState(instance)
		cpu_current := state.CPU.Usage
		dur := time.Now().Sub(checktime)
		elapsed_cpu := cpu_current-cpu
		cpuusage =elapsed_cpu * 100 / int64(dur)
	}else if instancetype == "vm"{
		state,_,_ := c.GetInstanceState(instance)
		sysInfo, err := pidusage.GetStat(int(state.Pid))
		if err != nil{
			fmt.Printf(err.Error())
			panic(err)
		}
		memusage = int64(sysInfo.Memory)
		cpu := state.CPU.Usage
		checktime := time.Now()
		sysInfo, err = pidusage.GetStat(int(state.Pid))
		if err != nil{
			fmt.Printf(err.Error())
			panic(err)
		}
		cpu_current := int64(sysInfo.CPU)
		dur := time.Now().Sub(checktime)
		elapsed_cpu := cpu_current-cpu
		cpuusage =elapsed_cpu * 100 / int64(dur)
	}
	return Sourceinfo{
		CpuUsage: cpuusage,
		MemUsage: memusage,
	}
}