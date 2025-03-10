//go:build windows

package wnlm

// NLMNetworkCategory represents the NLM_NETWORK_CATEGORY enumeration
// (a set of flags that specify the category type of a network).
//
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/ne-netlistmgr-nlm_network_category.
type NLMNetworkCategory byte

const (
	// NLMNetworkCategoryPublic represents the network category for public networks.
	NLMNetworkCategoryPublic = NLMNetworkCategory(0)
	// NLMNetworkCategoryPrivate represents the network category for private networks.
	NLMNetworkCategoryPrivate = NLMNetworkCategory(0x1)
	// NLMNetworkCategoryDomainAuthenticated represents the network category for domain authenticated networks.
	NLMNetworkCategoryDomainAuthenticated = NLMNetworkCategory(0x2)
)
