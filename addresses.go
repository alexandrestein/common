package common

import (
	"fmt"
	"net"
)

type (
	// Addr is an interface for net.Addr.
	Addr struct {
		MainAddr string
		Port     uint16
		Addrs    []string
	}
)

// NewAddr build a new Addr pointer with all addresses for every
// network interfaces with the first external IP as default or an error if any
func NewAddr(port uint16) (*Addr, error) {
	addrs, err := GetAddresses()
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		return nil, fmt.Errorf("no address found")
	}

	return &Addr{
		MainAddr: addrs[0],
		Port:     port,
		Addrs:    addrs,
	}, nil
}

// MustNewAddr does the same as NewAddr but returns an empty main address as "::"
func MustNewAddr(port uint16) *Addr {
	ret, err := NewAddr(port)
	if err != nil {
		return &Addr{
			MainAddr: "::",
			Port:     port,
			Addrs:    []string{},
		}
	}
	return ret
}

// AddAddr adds a new address to the list and returns it's position
func (a *Addr) AddAddr(new string) (addAtIndex int) {
	n := len(a.Addrs)
	a.Addrs = append(a.Addrs, new)
	return n
}

// AddAddrAndSwitch does the same as AddAddr but switch to the new address as main address
func (a *Addr) AddAddrAndSwitch(new string) {
	n := a.AddAddr(new)
	a.SwitchMain(n)
}

// SwitchMain switch the default address returned by Addr which is the port of the
// net.Addr interface
func (a *Addr) SwitchMain(i int) string {
	if i > len(a.Addrs)-1 {
		return ""
	}
	a.MainAddr = a.Addrs[i]
	return a.String()
}

// String returns the main address with the port as net.Addr expect it
func (a *Addr) String() string {
	return fmt.Sprintf("%s:%d", a.MainAddr, a.Port)
}

// Network returns "tcp", it's part of the net.Addr interface
func (a *Addr) Network() string {
	return "tcp"
}

// ForListenerBroadcast can be used for listeners to listen on any interfaces
// at the registered port
func (a *Addr) ForListenerBroadcast() string {
	return fmt.Sprintf(":%d", a.Port)
}

// IP returns the main address as a net.IP variable
func (a *Addr) IP() net.IP {
	return net.ParseIP(a.MainAddr)
}

// GetAddresses returns a slice of all detected external addresses or an error if any
func GetAddresses() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ret := []string{}

	for _, nic := range interfaces {
		var addrs []net.Addr
		addrs, err = nic.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ipAsString := addr.String()
			ip, _, err := net.ParseCIDR(ipAsString)
			if err != nil {
				continue
			}

			ipAsString = ip.String()
			ip2 := net.ParseIP(ipAsString)
			if to4 := ip2.To4(); to4 == nil {
				ipAsString = "[" + ipAsString + "]"
			}

			// If ip accessible from outside
			if ip.IsGlobalUnicast() {
				ret = append(ret, ipAsString)
			}
		}
	}

	return ret, nil
}
