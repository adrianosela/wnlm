//go:build windows

package wnlm

// NLMDomainType represents the NLM_DOMAIN_TYPE enumeration
// (a set of flags that specify the domain type of a network).
//
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/ne-netlistmgr-nlm_domain_type.
type NLMDomainType int32

const (
	// NLMDomainTypeNonDomainNetwork represents the domain type for non domain networks.
	NLMDomainTypeNonDomainNetwork = NLMDomainType(0)
	// NLMDomainTypeDomainNetwork represents the domain type for domain networks.
	NLMDomainTypeDomainNetwork = NLMDomainType(0x1)
	// NLMDomainTypeDomainAuthenticated represents the domain type for domain authenticated networks.
	NLMDomainTypeDomainAuthenticated = NLMDomainType(0x2)
)

var nlmDomainTypeToString = map[NLMDomainType]string{
	NLMDomainTypeNonDomainNetwork:    "None",
	NLMDomainTypeDomainNetwork:       "Domain",
	NLMDomainTypeDomainAuthenticated: "Domain Authenticated",
}

// String returns the string representation of the NLMDomainType.
func (t NLMDomainType) String() string {
	if str, ok := nlmDomainTypeToString[t]; ok {
		return str
	}
	return ""
}
