package dhcp

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lu1as/freeipa-dhcp/freeipa"
	ps "github.com/mitchellh/go-ps"
	log "github.com/sirupsen/logrus"
)

const ()

type DHCPConnector struct {
	ipa       *freeipa.FreeIPAClient
	hosts     map[string]*DHCPHostEntry
	update    time.Duration
	process   *os.Process
	connector SpecificDHCPConnector
}

type SpecificDHCPConnector interface {
	Notify(*os.Process) error
	Format(*DHCPHostEntry) string
}

func NewDHCPConnector(ipa *freeipa.FreeIPAClient, updateInterval time.Duration) *DHCPConnector {
	return &DHCPConnector{
		ipa:    ipa,
		hosts:  make(map[string]*DHCPHostEntry),
		update: updateInterval,
	}
}

func (c *DHCPConnector) Start(server string, zone string, filePath string, ttl string) {
	switch server {
	case "dnsmasq":
		c.connector = &DNSmasqConnector{}
		break
	default:
		log.Fatalf("%s dhcp server is not supported yet", server)
	}

	pid, err := c.getPid(server)
	if err != nil {
		log.Warnf(err.Error())
	} else if c.process, err = os.FindProcess(pid); err != nil {
		log.Warnf(err.Error())
	} else {
		log.Infof("found %s process with pid %d", server, pid)
	}

	log.Infof("checking every %s for new hosts", c.update)
	for {
		hosts, err := c.generate(zone, ttl)
		if err != nil {
			log.Fatal(err.Error())
		}

		if c.hasChanged(hosts) {
			c.hosts = hosts
			err = c.write(filePath)
			if err != nil {
				log.Fatal(err.Error())
			}

			log.Infof("hosts changed, reload %s", server)
			c.connector.Notify(c.process)
		}
		<-time.NewTimer(c.update).C
	}
}

func (c *DHCPConnector) generate(zone string, ttl string) (map[string]*DHCPHostEntry, error) {
	allHosts, err := c.ipa.GetHosts()
	if err != nil {
		return nil, err
	}

	hosts := make(map[string]*DHCPHostEntry)
	for _, i := range allHosts {
		if len(i.MAC) < 1 {
			log.Debugf("host %s has no mac address, skipping", i.Fqdn[0])
			continue
		}

		host := strings.SplitN(i.Fqdn[0], ".", 2)
		if len(host) < 2 {
			log.Debug("host has invalid fqdn")
			continue
		} else if host[1] != zone {
			log.Debugf("host %s is in another zone, skipping", i.Fqdn[0])
			continue
		}

		dns, err := c.ipa.GetDNSRecord(host[1], host[0])
		if err != nil {
			log.Debug(err.Error())
			continue
		} else if len(dns.ARecord) < 1 {
			log.Debugf("host %s has no A record, skipping", i.Fqdn[0])
			continue
		}

		he := &DHCPHostEntry{
			Name: host[0],
			IP:   dns.ARecord[0],
			MAC:  i.MAC[0],
			TTL:  ttl,
		}

		log.Debugf("add host: name: %s, ip: %s, mac: %s, ttl: %s", host[0], dns.ARecord[0], i.MAC[0], ttl)
		hosts[he.Name] = he
	}

	return hosts, nil
}

func (c *DHCPConnector) hasChanged(hosts map[string]*DHCPHostEntry) bool {
	if len(c.hosts) != len(hosts) {
		return true
	}

	for _, host := range c.hosts {
		if o, ok := hosts[host.Name]; !ok {
			return true
		} else if !host.equals(o) {
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

	return 0, fmt.Errorf("no dnsmasq process found")
}

func (c *DHCPConnector) write(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

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
