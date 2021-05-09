package main

import (
	"fmt"
	lxd "github.com/lxc/lxd/client"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func monitoer_srv(client lxd.InstanceServer) {
	router := mux.NewRouter()

	router.HandleFunc("/network/{name}/{nic}/{role}/{operation}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		nic := vars["nic"]
		prt := vars["operation"]
		role := vars["role"]
		status := true
		var statusinfo string
		if prt == "start" {
			go func() {
				for status {
					info := Getnetstat(name, nic, client)
					Netstats.WithLabelValues(role, name, nic, "BytesReceived").Set(float64(info.BytesReceived))
					Netstats.WithLabelValues(role, name, nic, "BytesSent").Set(float64(info.BytesSent))
					Netstats.WithLabelValues(role, name, nic, "PacketsReceived").Set(float64(info.PacketsReceived))
					Netstats.WithLabelValues(role, name, nic, "PacketsSent").Set(float64(info.PacketsSent))
					time.Sleep(time.Second * 5)
				}
			}()
			statusinfo = fmt.Sprintf("instance name:%s\ninterface:%s\nmonitor status: start", name, nic)
			_, _ = w.Write([]byte(statusinfo))
		} else if prt == "stop" {
			//todo: 停止监控的逻辑
			status = false
			statusinfo = fmt.Sprintf("instance name:%s\ninterface:%s\nmonitor status: stop", name, nic)
			_, _ = w.Write([]byte(statusinfo))
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(302)
	})
	router.HandleFunc("/instance/{name}/{role}/{operation}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(302)
		vars := mux.Vars(r)
		name := vars["name"]
		prt := vars["operation"]
		role := vars["role"]
		status := true
		var statusinfo string
		if prt == "start" {
			go func() {
				for status {
					info := Getsource(name, client)
					Source.WithLabelValues(role, name, "CPU").Set(float64(info.CpuUsage))
					Source.WithLabelValues(role, name, "MEM").Set(float64(info.MemUsage))
					time.Sleep(time.Second * 5)
				}
			}()
			statusinfo = fmt.Sprintf("instance name:%s\nmonitor status: start\n", name)
			_, _ = w.Write([]byte(statusinfo))
		} else if prt == "stop" {
			//todo: 停止监控的逻辑
			status = false
			statusinfo = fmt.Sprintf("instance name:%s\nmonitor status: stop\n", name)
			_, _ = w.Write([]byte(statusinfo))
		}
	})
	srv := &http.Server{
		Handler: router,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
