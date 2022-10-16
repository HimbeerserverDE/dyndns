// dyndns is a dual-stack DynDNS client that uses the INWX JSON-RPC API.
package main

import (
	"flag"
	"net"
)

func main() {
	const usage = "override the default config path"
	configFile := flag.String("config", "/etc/dyndns.conf", usage)

	conf := &config{}
	if err := conf.parse(*configFile); err != nil {
		logger.Fatal(err)
	}

	update4 := make(chan net.IPAddr)
	update6 := make(chan net.IPNet)

	go monitor(conf, update4, update6)

	for {
		select {
		case newAddr := <-update4:
			// TODO
		case newPrefix := <-update6:
			// TODO
		}
	}
}
