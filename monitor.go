package main

import (
	"net"
	"time"
)

func monitor4(conf *config, update4 chan<- net.IPAddr) {
	refresh := time.NewTicker(conf.Interval)
	for range refresh.C {
		link, err := net.InterfaceByName(conf.Link4)
		if err != nil {
			logger.Printf("can't find interface %s: %w", conf.Link4, err)
			continue
		}

		addrs, err := link.Addrs()
		if err != nil {
			const format = "can't read IPv4 addresses from %s: %w"
			logger.Printf(format, conf.Link4, err)
			continue
		}

		var addr4 *net.IPAddr
		for _, netAddr := range addrs {
			addr := netAddr.(*net.IPAddr)

			if addr.To4() != nil && !addr.IP.IsPrivate() {
				addr4 = addr
				break
			}
		}
	}
}

func monitor6(conf *config, update6 chan<- net.IPNet) {
	refresh := time.NewTicker(conf.Interval)
	for range refresh.C {
		link, err := net.InterfaceByName(conf.Link6)
		if err != nil {
			logger.Printf("can't find interface %s: %w", conf.Link6, err)
			continue
		}

		addrs, err := link.Addrs()
		if err != nil {
			const format = "can't read IPv6 addresses from %s: %w"
			logger.Printf(format, conf.Link6, err)
			continue
		}

		var prefix6 *net.IPNet
		for _, netAddr := range addrs {
			addr := netAddr.(*net.IPAddr)

			if addr.To4() == nil && addr.IP.IsGlobalUnicast() {
				cidr := net.CIDRMask(conf.PrefixLen, 128)
				prefix6 = addr.IP.Mask(cidr)

				break
			}
		}
	}
}
