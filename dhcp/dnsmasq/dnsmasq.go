package dnsmasq

import (
	"fmt"
	"os"
	"syscall"

	"github.com/lu1as/freeipa-dhcp/dhcp"
)

type DNSmasqConnector struct{}

func (c *DNSmasqConnector) Reload(p *os.Process) error {
	if p != nil {
		if err := p.Signal(syscall.SIGHUP); err != nil {
			return err
		}
	}
	return nil
}

func (c *DNSmasqConnector) Format(e *dhcp.DHCPHostEntry) string {
	return fmt.Sprintf("%s,%s,%s,infinite\n", e.MAC, e.Name, e.IP)
}
