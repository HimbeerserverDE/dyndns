// dyndns is a dual-stack DynDNS client that uses the INWX JSON-RPC API.
package main

import (
	"flag"
	"fmt"
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
			if err := nsUpdate4(conf, newAddr); err != nil {
				logger.Print(err)
			}
		case newPrefix := <-update6:
			if err := nsUpdate6(conf, newPrefix); err != nil {
				logger.Print(err)
			}
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

func nsUpdate4(conf *config, addr *net.IPAddr) error {
	clt, err := inwx.Login(inwx.Production, conf.User, conf.Passwd)
	if err != nil {
		return err
	}
	defer clt.Close()

	return updateRecords(clt, conf.Records4, addr.String())
}

func nsUpdate6(conf *config, prefix *net.IPNet) error {
	return fmt.Errorf("updating IPv6 records is not yet implemented")
}
