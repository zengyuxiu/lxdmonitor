package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func monitoer_srv() {
	router := mux.NewRouter()

	router.HandleFunc("/network/{name}/{nic}/{operation}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		nic := vars["nic"]
		prt := vars["operation"]
		status := true
		var statusinfo string
		if prt == "start"{
			go func() {
				for  status {
					info := Getnetstat(name,nic)
					Netstats.WithLabelValues(name,nic,"BytesReceived").Observe(float64(info.BytesReceived))
					Netstats.WithLabelValues(name,nic,"BytesSent").Observe(float64(info.BytesSent))
					Netstats.WithLabelValues(name,nic,"PacketsReceived").Observe(float64(info.PacketsReceived))
					Netstats.WithLabelValues(name,nic,"PacketsSent").Observe(float64(info.PacketsSent))
					time.Sleep(time.Second*5)
				}
			}()
			statusinfo = fmt.Sprintf("instance name:%s\ninterface:%s\nmonitor status: start",name,nic)
			_, _ = w.Write([]byte(statusinfo))
		}else if prt == "stop"{
			//todo: 停止监控的逻辑
			status = false
			statusinfo = fmt.Sprintf("instance name:%s\ninterface:%s\nmonitor status: stop",name,nic)
			_, _ = w.Write([]byte(statusinfo))
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(302)
	})
	router.HandleFunc("/instance/{name}/{type}/{operation}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(302)
		vars := mux.Vars(r)
		name := vars["name"]
		types := vars["type"]
		prt := vars["operation"]
		status := true
		var statusinfo string
		if prt == "start"{
			go func() {
				for  status {
					info := Getsource(name,types)
					Source.WithLabelValues(types,name,"CPU").Set(float64(info.CpuUsage))
					Source.WithLabelValues(types,name,"MEM").Set(float64(info.MemUsage))
					time.Sleep(time.Second*5)
				}
			}()
			statusinfo = fmt.Sprintf("instance name:%s\ninstance type:%s\nmonitor status: start",name,types)
			_, _ = w.Write([]byte(statusinfo))
		}else if prt == "stop"{
			//todo: 停止监控的逻辑
			status = false
			statusinfo = fmt.Sprintf("instance name:%s\ninstance type:%s\nmonitor status: stop",name,types)
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
