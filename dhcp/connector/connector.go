package connector

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lu1as/freeipa-dhcp/dhcp"
	"github.com/lu1as/freeipa-dhcp/dhcp/dhcpd"
	"github.com/lu1as/freeipa-dhcp/dhcp/dnsmasq"
	"github.com/lu1as/freeipa-dhcp/freeipa"
	ps "github.com/mitchellh/go-ps"
	log "github.com/sirupsen/logrus"
)

type DHCPConnector struct {
	ipa       *freeipa.FreeIPAClient
	hosts     map[string]*dhcp.DHCPHostEntry
	update    time.Duration
	process   *os.Process
	connector dhcp.DHCPConnector
}

func NewDHCPConnector(ipa *freeipa.FreeIPAClient, updateInterval time.Duration) *DHCPConnector {
	return &DHCPConnector{
		ipa:    ipa,
		hosts:  make(map[string]*dhcp.DHCPHostEntry),
		update: updateInterval,
	}
}

func (c *DHCPConnector) Start(server string, zone string, filePath string) {
	switch server {
	case "dnsmasq":
		c.connector = &dnsmasq.DNSmasqConnector{}
		break
	case "dhcpd":
		c.connector = &dhcpd.DHCPdConnector{}
	default:
		log.Fatalf("%s dhcp server is not supported yet", server)
	}

	pid, err := c.getPid(server)
	if err != nil {
		log.Warnf("DHCP server reload disabled: %s", err.Error())
	} else if c.process, err = os.FindProcess(pid); err != nil {
		log.Warnf("DHCP server reload disabled: %s", err.Error())
	} else {
		log.Infof("Found %s process with pid %d", server, pid)
	}

	log.Infof("Checking every %s for new hosts", c.update)
	for {
		hosts, err := c.generate(zone)
		if err != nil {
			log.Fatal(err.Error())
		}

		if c.hasChanged(hosts) {
			c.hosts = hosts
			err = c.write(filePath)
			if err != nil {
				log.Fatal(err.Error())
			}

			log.Infof("Hosts changed, reload %s", server)
			c.connector.Reload(c.process)
		}
		<-time.NewTimer(c.update).C
	}
}

func (c *DHCPConnector) generate(zone string) (map[string]*dhcp.DHCPHostEntry, error) {
	allHosts, err := c.ipa.GetHosts()
	if err != nil {
		return nil, err
	}

	hosts := make(map[string]*dhcp.DHCPHostEntry)
	for _, i := range allHosts {
		if len(i.MAC) < 1 {
			log.Debugf("Host %s has no mac address, skipping", i.Fqdn[0])
			continue
		}

		host := strings.SplitN(i.Fqdn[0], ".", 2)
		if len(host) < 2 {
			log.Debug("Host has invalid fqdn")
			continue
		} else if host[1] != zone {
			log.Debugf("Host %s is in another zone, skipping", i.Fqdn[0])
			continue
		}

		dns, err := c.ipa.GetDNSRecord(host[1], host[0])
		if err != nil {
			log.Debug(err.Error())
			continue
		} else if len(dns.ARecord) < 1 {
			log.Debugf("Host %s has no A record, skipping", i.Fqdn[0])
			continue
		}

		he := &dhcp.DHCPHostEntry{
			Name: host[0],
			IP:   dns.ARecord[0],
			MAC:  i.MAC[0],
		}

		log.Debugf("Add host: name: %s, ip: %s, mac: %s", host[0], dns.ARecord[0], i.MAC[0])
		hosts[he.Name] = he
	}

	return hosts, nil
}

func (c *DHCPConnector) hasChanged(hosts map[string]*dhcp.DHCPHostEntry) bool {
	if len(c.hosts) != len(hosts) {
		return true
	}

	for _, host := range c.hosts {
		if o, ok := hosts[host.Name]; !ok {
			return true
		} else if !host.Equals(o) {
			return true
		}
	}
	return false
}

func (c *DHCPConnector) getPid(name string) (int, error) {
	all, err := ps.Processes()
	if err != nil {
		return 0, err
	}

	for _, p := range all {
		if strings.HasPrefix(p.Executable(), name) {
			return p.Pid(), nil
		}
	}

	return 0, fmt.Errorf("No %s process found", name)
}

func (c *DHCPConnector) write(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return err
	}

	for _, host := range c.hosts {
		if _, err := f.WriteString(c.connector.Format(host)); err != nil {
			return err
		}
	}

	if err := f.Sync(); err != nil {
		return err
	}
	f.Close()

	return nil
}
