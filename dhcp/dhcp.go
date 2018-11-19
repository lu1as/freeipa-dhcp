package dhcp

import (
	"os"
)

type DHCPConnector interface {
	Reload(*os.Process) error
	Format(*DHCPHostEntry) string
}

type DHCPHostEntry struct {
	MAC  string
	Name string
	IP   string
}

func (e *DHCPHostEntry) Equals(o *DHCPHostEntry) bool {
	if e.Name == o.Name &&
		e.IP == o.IP &&
		e.MAC == o.MAC {
		return true
	}
	return false
}
