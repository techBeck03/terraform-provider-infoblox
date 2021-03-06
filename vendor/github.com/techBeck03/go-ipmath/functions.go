package ipmath

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IP is the base building block
type IP struct {
	Address net.IP
	Network *net.IPNet
}

// NewIP creates new ipmath base object
func NewIP(cidrAddr string) (newIP IP, err error) {
	ip, net, err := net.ParseCIDR(cidrAddr)
	if err != nil {
		return newIP, err
	}
	return IP{
		Address: ip,
		Network: net,
	}, nil
}

// Add increments IP address by supplied increment value
func (i *IP) Add(incr int) error {
	ip := i.Address.To4()
	if ip == nil {
		return fmt.Errorf("%s is not a valid IP address", ip.String())
	}

	incBit := uint32(ip[3]) | uint32(ip[2])<<8 | uint32(ip[1])<<16 | uint32(ip[0])<<24
	incBit += uint32(incr)
	newIP := net.ParseIP(strings.Join([]string{fmt.Sprint(incBit & 0xff000000 >> 24), fmt.Sprint(incBit & 0x00ff0000 >> 16), fmt.Sprint(incBit & 0x0000ff00 >> 8), fmt.Sprint(incBit & 0x000000ff)}, "."))

	if newIP == nil {
		return fmt.Errorf("%s cannot be incremented by %v", i.Address.String(), incr)
	}
	if i.Network != nil {
		if !i.Network.Contains(newIP) {
			return fmt.Errorf("%s is not in CIDR network %s", newIP.String(), i.Network.String())
		}
	}
	if bytes.Compare(i.Address, newIP) > 0 {
		return fmt.Errorf("Adding %v to %s causes address wrap", incr, i.Address.String())
	}
	i.Address = newIP
	return nil
}

// Inc increments IP address by one
func (i *IP) Inc() error {
	return i.Add(1)
}

// Subtract decrements IP address by supplied increment value
func (i *IP) Subtract(incr int) error {
	ip := i.Address.To4()
	if ip == nil {
		return fmt.Errorf("%s is not a valid IP address", ip.String())
	}

	incBit := uint32(ip[3]) | uint32(ip[2])<<8 | uint32(ip[1])<<16 | uint32(ip[0])<<24
	incBit -= uint32(incr)
	newIP := net.ParseIP(strings.Join([]string{fmt.Sprint(incBit & 0xff000000 >> 24), fmt.Sprint(incBit & 0x00ff0000 >> 16), fmt.Sprint(incBit & 0x0000ff00 >> 8), fmt.Sprint(incBit & 0x000000ff)}, "."))
	if newIP == nil {
		return fmt.Errorf("%s cannot be incremented by %v", i.Address.String(), incr)
	}
	if i.Network != nil {
		if !i.Network.Contains(newIP) {
			return fmt.Errorf("%s is not in CIDR network %s", newIP.String(), i.Network.String())
		}
	}
	if bytes.Compare(newIP, i.Address) > 0 {
		return fmt.Errorf("Subtracting %v to %s causes address wrap", incr, i.Address.String())
	}
	i.Address = newIP
	return nil
}

// Difference find the signed difference between two IPs
func (i *IP) Difference(ip net.IP) int {
	ipA := i.Address.To4()
	ipABit := uint32(ipA[3]) | uint32(ipA[2])<<8 | uint32(ipA[1])<<16 | uint32(ipA[0])<<24
	ipB := ip.To4()
	ipBBit := uint32(ipB[3]) | uint32(ipB[2])<<8 | uint32(ipB[1])<<16 | uint32(ipB[0])<<24

	return int(ipBBit) - int(ipABit)
}

// Dec decrements IP address by one
func (i *IP) Dec() error {
	return i.Subtract(1)
}

// EQ checks if ip a is greater than ip b
func (i *IP) EQ(ip net.IP) bool {
	return bytes.Compare(i.Address, ip) == 0
}

// GT checks if ip a is greater than ip b
func (i *IP) GT(ip net.IP) bool {
	return bytes.Compare(i.Address, ip) > 0
}

// GTE checks if ip a is greater than or equal to ip b
func (i *IP) GTE(ip net.IP) bool {
	return bytes.Compare(i.Address, ip) >= 0
}

// LT checks if ip a is less than ip b
func (i *IP) LT(ip net.IP) bool {
	return bytes.Compare(i.Address, ip) < 0
}

// LTE checks if ip a is less than or equal to ip b
func (i *IP) LTE(ip net.IP) bool {
	return bytes.Compare(i.Address, ip) <= 0
}

// ToIPString prints ipmath object host ip address
func (i *IP) ToIPString() string {
	return i.Address.String()
}

// ToCIDRString prints ipmath object in CIDR format
func (i *IP) ToCIDRString() (string, error) {
	if i.Network == nil {
		return i.Address.String(), fmt.Errorf("Unable to create cidr string because `Network` is undefined")
	}
	prefix, _ := i.Network.Mask.Size()
	return fmt.Sprintf("%s/%s", i.Address.String(), strconv.Itoa(prefix)), nil
}

// Clone clones the ipmath base object
func (i *IP) Clone() IP {
	if i.Network != nil {
		cidrString, _ := i.ToCIDRString()
		newIP, _ := NewIP(cidrString)
		return newIP
	}
	return IP{
		Address: net.ParseIP(i.Address.String()),
	}
}
