package dhcpd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lu1as/freeipa-dhcp/dhcp"
)

type DHCPdConnector struct{}

func (c *DHCPdConnector) Reload(p *os.Process) error {
	if p != nil {
		cmd := exec.Command("systemctl", "restart", "dhcpd")
		return cmd.Run()
	}
	return nil
}

func (c *DHCPdConnector) Format(e *dhcp.DHCPHostEntry) string {
	return fmt.Sprintf("host %s { hardware ethernet %s; fixed-address %s; }\n", e.Name, e.MAC, e.IP)
}
