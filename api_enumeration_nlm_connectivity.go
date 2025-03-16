//go:build windows

package wnlm

import (
	"sort"
	"strings"

	"github.com/adrianosela/wnlm/pkg/bits"
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

// IsDisconnected returns true if the NLMConnectivity has the disconnected flag set.
func (c NLMConnectivity) IsDisconnected() bool {
	return c == 0
}

// IsIPv4NoTraffic returns true if the NLMConnectivity has IPv4NoTraffic flag set.
func (c NLMConnectivity) IsIPv4NoTraffic() bool {
	return bits.AreSet(c, NLMConnectivityIPv4NoTraffic)
}

// IsIPv6NoTraffic returns true if the NLMConnectivity has the IPv6NoTraffic flag set.
func (c NLMConnectivity) IsIPv6NoTraffic() bool {
	return bits.AreSet(c, NLMConnectivityIPv6NoTraffic)
}

// IsIPv4Subnet returns true if the NLMConnectivity has the IPv4Subnet flag set.
func (c NLMConnectivity) IsIPv4Subnet() bool {
	return bits.AreSet(c, NLMConnectivityIPv4Subnet)
}

// IsIPv4LocalNetwork returns true if the NLMConnectivity has the IPv4LocalNetwork flag set.
func (c NLMConnectivity) IsIPv4LocalNetwork() bool {
	return bits.AreSet(c, NLMConnectivityIPv4LocalNetwork)
}

// IsIPv4Internet returns true if the NLMConnectivity has the IPv4Internet flag set.
func (c NLMConnectivity) IsIPv4Internet() bool {
	return bits.AreSet(c, NLMConnectivityIPv4Internet)
}

// IsIPv6Subnet returns true if the NLMConnectivity has the IPv6Subnet flag set.
func (c NLMConnectivity) IsIPv6Subnet() bool {
	return bits.AreSet(c, NLMConnectivityIPv6Subnet)
}

// IsIPv6LocalNetwork returns true if the NLMConnectivity has the IPv6LocalNetwork flag set.
func (c NLMConnectivity) IsIPv6LocalNetwork() bool {
	return bits.AreSet(c, NLMConnectivityIPv6LocalNetwork)
}

// IsIPv6Internet returns true if the NLMConnectivity has the IPv6Internet flag set.
func (c NLMConnectivity) IsIPv6Internet() bool {
	return bits.AreSet(c, NLMConnectivityIPv6Internet)
}

// Strings provides a string representation of the NLMConnectivity status.
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
