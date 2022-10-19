// dyndns is a dual-stack DynDNS client that uses the INWX JSON-RPC API.
package main

import (
	"flag"
	"net"

	"github.com/HimbeerserverDE/inwx"
)

func main() {
	const usage = "override the default config path"
	configFile := flag.String("config", "/etc/dyndns.conf", usage)

	conf := &config{}
	if err := conf.parse(*configFile); err != nil {
		logger.Fatal(err)
	}

	update4 := make(chan *net.IPAddr)
	update6 := make(chan *net.IPNet)

	go monitor4(conf, update4)
	go monitor6(conf, update6)

	for {
		select {
		case newAddr := <-update4:
			// TODO
		case newPrefix := <-update6:
			// TODO
		}
	}
}

func updateRecords(c *inwx.Client, ids []int, content string) error {
	_, err := c.Call(&inwx.NSUpdateRecordsCall{
		IDs: ids,
		RecordInfo: inwx.RecordInfo{
			Content: content,
		},
	})

	return err
}
