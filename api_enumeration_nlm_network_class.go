//go:build windows

package wnlm

// NLMNetworkClass represents the NLM_NETWORK_CLASS enum (a set
// of flags that that specify if a network has been identified).
//
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/ne-netlistmgr-nlm_network_class.
type NLMNetworkClass int32

const (
	// NLMNetworkClassIdentified represents the network class for identifying networks.
	NLMNetworkClassIdentifying = NLMNetworkClass(0x1)
	// NLMNetworkClassIdentified represents the network class for identified networks.
	NLMNetworkClassIdentified = NLMNetworkClass(0x2)
	// NLMNetworkClassUnidentified represents the network class for unidentified networks.
	NLMNetworkClassUnidentified = NLMNetworkClass(0x3)
)

var nlmNetworkClassToString = map[NLMNetworkClass]string{
	NLMNetworkClassIdentifying:  "Identifying",
	NLMNetworkClassIdentified:   "Identified",
	NLMNetworkClassUnidentified: "Unidentified",
}

// String returns the string representation of the NLMNetworkClass.
func (c NLMNetworkClass) String() string {
	if str, ok := nlmNetworkClassToString[c]; ok {
		return str
	}
	return ""
}
