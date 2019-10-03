package publisher

import (
	"fmt"
	"net"
	"sync"

	"github.com/celebdor/zeroconf"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Publish(ip net.IP, iface net.Interface, service Service, shutdown chan struct{}, waitGroup *sync.WaitGroup) (err error) {
	defer waitGroup.Done()
	svcEntry := zeroconf.NewServiceEntry(service.Name, service.SvcType, service.Domain)
	svcEntry.Port = service.Port
	if ip.To4() != nil {
		svcEntry.AddrIPv4 = append(svcEntry.AddrIPv4, ip)
	} else {
		svcEntry.AddrIPv6 = append(svcEntry.AddrIPv6, ip)
	}
	svcEntry.HostName = service.HostName
	log.WithFields(logrus.Fields{
		"name": svcEntry.Instance,
	}).Info("Zeroconf registering service")
	s, err := zeroconf.RegisterSvcEntry(svcEntry, []net.Interface{iface})
	if err != nil {
		log.Error("Failed to create zeroconf Server", err)
		return err
	}
	defer s.Shutdown()
	log.WithFields(logrus.Fields{
		"name": svcEntry.Instance,
		"ttl":  service.TTL,
	}).Info("Zeroconf setting service ttl")
	s.TTL(service.TTL)

	select {
	case <-shutdown:
		log.WithFields(logrus.Fields{
			"name": svcEntry.Instance,
		}).Info("Gracefully shutting down service")
	}

	return nil
}

func FindIface(ip net.IP) (iface net.Interface, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("[ERR] mdns-publish: Failed retrieving system network interfaces %v.", err)
		return iface, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Printf("[ERR] mdns-publish: Failed retrieving network addresses for interface %s: %v.", i.Name, err)
		}
		for _, addr := range addrs {
			var currIP net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				currIP = v.IP
			case *net.IPAddr:
				currIP = v.IP
			}
			if currIP == nil {
				continue
			}
			if currIP.Equal(ip) {
				iface = i
				return iface, nil
			}
		}
	}
	return iface, fmt.Errorf("Couldn't find interface with IP address %s", ip)
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}
