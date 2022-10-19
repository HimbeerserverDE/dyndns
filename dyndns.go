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
	clt, err := inwx.Login(inwx.Production, conf.User, conf.Passwd)
	if err != nil {
		return err
	}
	defer clt.Close()

	for _, id := range conf.Records6 {
		resp, err := clt.Call(&inwx.NSRecordInfoCall{
			RecordID: id,
			Type:     "AAAA",
		})
		if err != nil {
			logger.Printf("can't get info, skipping IPv6 record %d: %v", id, err)
			continue
		}

		record := &inwx.NSRecordInfoResponse{}
		if err := resp.Into(record); err != nil {
			logger.Printf("IPv6 record response is invalid: %v", err)
			continue
		}

		if len(record.Record) != 1 {
			logger.Printf("invalid number of IPv6 records: %d != 1", len(record.Record))
			continue
		}

		addr := net.ParseIP(record.Record[0].Content)
		if addr == nil {
			logger.Printf("invalid IPv6 record %d: %s", id, record.Record[0].Content)
			continue
		}

		ifidMask := net.CIDRMask(conf.PrefixLen, 128)
		for k, v := range ifidMask {
			ifidMask[k] = ^v
		}

		ifid := addr.Mask(ifidMask)

		newAddr := prefix.IP
		for k, v := range newAddr {
			newAddr[k] = v | ifid[k]
		}

		if err := updateRecords(clt, []int{id}, newAddr.String()); err != nil {
			logger.Printf("can't update, skipping IPv6 record %d: %v", id, err)
			continue
		}
	}

	return nil
}
