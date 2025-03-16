//go:build windows

package wnlm

import (
	"sort"
	"strings"

	"github.com/adrianosela/wnlm/pkg/bits"
)

// NLMInternetConnectivity represents the NLM_INTERNET_CONNECTIVITY enum (a set
// of flags that provide additional data for IPv4 or IPv6 network connectivity).
//
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/ne-netlistmgr-nlm_internet_connectivity.
type NLMInternetConnectivity int32

const (
	// NLMInternetConnectivityWebHijack represents the web hijack Internet connectivity.
	NLMInternetConnectivityWebHijack = NLMInternetConnectivity(0x1)
	// NLMInternetConnectivityProxied represents the proxied Internet connectivity.
	NLMInternetConnectivityProxied = NLMInternetConnectivity(0x2)
	// NLMInternetConnectivityCorporate represents the corporate Internet connectivity.
	NLMInternetConnectivityCorporate = NLMInternetConnectivity(0x4)
)

// IsWebHijack returns true if the NLMInternetConnectivity has the WebHijack flag set.
func (c NLMInternetConnectivity) IsWebHijack() bool {
	return bits.AreSet(c, NLMInternetConnectivityWebHijack)
}

// IsProxied returns true if the NLMInternetConnectivity has the Proxied flag set.
func (c NLMInternetConnectivity) IsProxied() bool {
	return bits.AreSet(c, NLMInternetConnectivityProxied)
}

// IsCorporate returns true if the NLMInternetConnectivity has the Corporate flag set.
func (c NLMInternetConnectivity) IsCorporate() bool {
	return bits.AreSet(c, NLMInternetConnectivityCorporate)
}

// Strings provides a string representation of the Internet connectivity.
func (c NLMInternetConnectivity) String() string {
	flags := []string{}
	for flag, check := range map[string]bool{
		"WebHijack": c.IsWebHijack(),
		"Proxied":   c.IsProxied(),
		"Corporate": c.IsCorporate(),
	} {
		if check {
			flags = append(flags, flag)
		}
	}
	sort.Strings(flags)
	return strings.Join(flags, ", ")
}
