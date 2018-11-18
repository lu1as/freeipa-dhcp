package dhcp

import (
	"fmt"
	"os"
	"syscall"
)

type DNSmasqConnector struct{}

func (c *DNSmasqConnector) Notify(p *os.Process) error {
	if p != nil {
		if err := p.Signal(syscall.SIGHUP); err != nil {
			return err
		}
	}
	return nil
}

func (c *DNSmasqConnector) Format(e *DHCPHostEntry) string {
	return fmt.Sprintf("%s,%s,%s,%s\n", e.MAC, e.Name, e.IP, e.TTL)
}
