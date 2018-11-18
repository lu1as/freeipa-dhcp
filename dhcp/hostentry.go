package dhcp

type DHCPHostEntry struct {
	MAC  string
	Name string
	IP   string
	TTL  string
}

func (e *DHCPHostEntry) equals(o *DHCPHostEntry) bool {
	if e.Name == o.Name &&
		e.IP == o.IP &&
		e.MAC == o.MAC &&
		e.TTL == o.TTL {
		return true
	}
	return false
}
