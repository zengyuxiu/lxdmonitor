package main

import (
	"flag"
	"github.com/lxc/lxd/shared/logger"
)

func main() {
	var host string
	flag.StringVar(&host, "host", "127.0.0.1", "host ip address")
	go prometheus_srv()
	d, err := InitLxdInstanceServer(host)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	monitoer_srv(*d)
}
