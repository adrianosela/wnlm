//go:build windows

package wnlm

import (
	"sort"
	"strings"
)

// NLMConnectivity represents the NLM_CONNECTIVITY enumeration (a set of flags that
// provide notification whenever connectivity related parameters have changed).
//
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/ne-netlistmgr-nlm_connectivity
type NLMConnectivity int32

const (
	// NLMConnectivityDisconnected represents the connectivity for disconnected networks.
	NLMConnectivityDisconnected = NLMConnectivity(0x0000)
	// NLMConnectivityIPv4NoTraffic represents the connectivity for IPv4 no-traffic networks.
	NLMConnectivityIPv4NoTraffic = NLMConnectivity(0x0001)
	// NLMConnectivityIPv6NoTraffic represents the connectivity for IPv6 no-traffic networks.
	NLMConnectivityIPv6NoTraffic = NLMConnectivity(0x0002)
	// NLMConnectivityIPv4Subnet represents the connectivity for IPv4 subnet networks.
	NLMConnectivityIPv4Subnet = NLMConnectivity(0x0010)
	// NLMConnectivityIPv4LocalNetwork represents the connectivity for IPv4 local network networks.
	NLMConnectivityIPv4LocalNetwork = NLMConnectivity(0x0020)
	// NLMConnectivityIPv4Internet represents the connectivity for IPv4 Internet networks.
	NLMConnectivityIPv4Internet = NLMConnectivity(0x0040)
	// NLMConnectivityIPv6Subnet represents the connectivity for IPv6 subnet networks.
	NLMConnectivityIPv6Subnet = NLMConnectivity(0x0100)
	// NLMConnectivityIPv6LocalNetwork represents the connectivity for IPv6 local network networks.
	NLMConnectivityIPv6LocalNetwork = NLMConnectivity(0x0200)
	// NLMConnectivityIPv6Internet represents the connectivity for IPv6 Internet networks.
	NLMConnectivityIPv6Internet = NLMConnectivity(0x0400)
)

func (c NLMConnectivity) IsDisconnected() bool     { return c == 0 }
func (c NLMConnectivity) IsIPv4NoTraffic() bool    { return c.is(NLMConnectivityIPv4NoTraffic) }
func (c NLMConnectivity) IsIPv6NoTraffic() bool    { return c.is(NLMConnectivityIPv6NoTraffic) }
func (c NLMConnectivity) IsIPv4Subnet() bool       { return c.is(NLMConnectivityIPv4Subnet) }
func (c NLMConnectivity) IsIPv4LocalNetwork() bool { return c.is(NLMConnectivityIPv4LocalNetwork) }
func (c NLMConnectivity) IsIPv4Internet() bool     { return c.is(NLMConnectivityIPv4Internet) }
func (c NLMConnectivity) IsIPv6Subnet() bool       { return c.is(NLMConnectivityIPv6Subnet) }
func (c NLMConnectivity) IsIPv6LocalNetwork() bool { return c.is(NLMConnectivityIPv6LocalNetwork) }
func (c NLMConnectivity) IsIPv6Internet() bool     { return c.is(NLMConnectivityIPv6Internet) }

// is checks if the connectivity includes the specific flag
func (c NLMConnectivity) is(flag NLMConnectivity) bool { return c&flag == flag }

// Strings provides a string representation of the connectivity status
func (c NLMConnectivity) String() string {
	if c.IsDisconnected() {
		return "Disconnected"
	}

	flags := []string{}
	for flag, check := range map[string]bool{
		"IPv4NoTraffic":    c.IsIPv4NoTraffic(),
		"IPv6NoTraffic":    c.IsIPv6NoTraffic(),
		"IPv4Subnet":       c.IsIPv4Subnet(),
		"IPv4LocalNetwork": c.IsIPv4LocalNetwork(),
		"IPv4Internet":     c.IsIPv4Internet(),
		"IPv6Subnet":       c.IsIPv6Subnet(),
		"IPv6LocalNetwork": c.IsIPv6LocalNetwork(),
		"IPv6Internet":     c.IsIPv6Internet(),
	} {
		if check {
			flags = append(flags, flag)
		}
	}
	sort.Strings(flags)
	return strings.Join(flags, ", ")
}
