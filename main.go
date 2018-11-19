package main

import (
	"time"

	"github.com/lu1as/freeipa-dhcp/dhcp/connector"
	"github.com/lu1as/freeipa-dhcp/freeipa"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

var (
	debug         bool
	update        time.Duration
	ipaServer     string
	ipaUser       string
	ipaPassword   string
	ipaInsecure   bool
	dhcpZone      string
	dhcpServer    string
	dhcpTTL       string
	dhcpHostsFile string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "set log level to debug")
	flag.DurationVar(&update, "updateInterval", time.Minute, "update interval")

	// freeipa config
	flag.StringVar(&ipaServer, "ipaServer", "https://master.ipa.example.de", "freeipa server address")
	flag.StringVar(&ipaUser, "ipaUser", "admin", "freeipa user name")
	flag.StringVar(&ipaPassword, "ipaPassword", "secret", "freeipa user password")
	flag.BoolVar(&ipaInsecure, "ipaInsecure", false, "allow invalid ssl certificate")

	// dhcp config
	flag.StringVar(&dhcpServer, "dhcpServer", "dnsmasq", "dhcp server name")
	flag.StringVar(&dhcpZone, "dhcpZone", "example.de", "dns zone")
	flag.StringVar(&dhcpTTL, "dhcpTTL", "infinite", "dhcp lease time to live")
	flag.StringVar(&dhcpHostsFile, "dhcpHostsFile", "zone.hosts", "hosts output file path")
}

func main() {
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	ipa := freeipa.NewFreeIPAClient(ipaServer)
	if ipaInsecure {
		ipa.AllowInsecure()
	}
	if err := ipa.Login(ipaUser, ipaPassword); err != nil {
		log.Fatal(err.Error())
	}

	d := connector.NewDHCPConnector(ipa, update)
	d.Start(dhcpServer, dhcpZone, dhcpHostsFile, dhcpTTL)
}
