// dyndns is a dual-stack DynDNS client that uses the INWX JSON-RPC API.
package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	const usage = "override the default config path"
	configFile := flag.String("config", "/etc/dyndns.conf", usage)

	config := &config{}
	if err := config.parse(*configFile); err != nil {
		log.Fatal(err)
	}

	update4 := make(chan net.IPAddr)
	update6 := make(chan net.IPNet)

	go monitor4(update4)
	go monitor6(update6)

	for {
		select {
		case newAddr := <-update4:
			// TODO
		case newPrefix := <-update6:
			// TODO
		}
	}
}
