package main

import (
	"net"
	"time"
)

func monitor4(conf *config, update4 chan<- *net.IPAddr) {
	refresh := time.NewTicker(conf.Interval)

	var prevAddr4 *net.IPAddr
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

			if addr.IP.To4() != nil && !addr.IP.IsPrivate() {
				addr4 = addr
				break
			}
		}

		if !addr4.IP.Equal(prevAddr4.IP) {
			logger.Println("detected new IPv4 address:", addr4)
			update4 <- addr4

			prevAddr4 = addr4
		}
	}
}

func monitor6(conf *config, update6 chan<- *net.IPNet) {
	refresh := time.NewTicker(conf.Interval)

	var prevPrefix6 *net.IPNet
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

			if addr.IP.To4() == nil && addr.IP.IsGlobalUnicast() {
				cidr := net.CIDRMask(conf.PrefixLen, 128)
				prefix6 = addr.IP.Mask(cidr)

				break
			}
		}

		if !prefix6.Equal(prevPrefix6) {
			logger.Println("detected new IPv6 prefix:", prefix6)
			update6 <- prefix6

			prevPrefix6 = prefix6
		}
	}
}
