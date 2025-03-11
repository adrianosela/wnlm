//go:build windows

package wnlm

import (
	"fmt"

	"github.com/go-ole/go-ole"
)

// INetwork represents the Windows INetwork type as defined in
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/nn-netlistmgr-inetwork.
//
// The Windows Global Unique Identifier (GUID) for this interface is DCB00002-570F-4A9B-8D69-199FDBA5723B.
type INetwork interface {
	// TODO: implement:
	// - properties: get_IsConnected, getIsConnectedToInternet
	// - methods: GetNetworkId, GetTimeCreatedAndConnected

	GetCategory() (NLMNetworkCategory, error)
	GetConnectivity() (NLMConnectivity, error)
	GetDescription() (string, error)
	GetDomainType() (NLMDomainType, error)
	GetName() (string, error)
	GetNetworkConnections() (IEnumNetworkConnections, error)
	SetCategory(NLMNetworkCategory) error
	SetDescription(string) error
	SetName(string) error

	Release()
}

// iNetwork is the default implementation of INetwork.
type iNetwork struct {
	idispatch *ole.IDispatch
}

// GetCategory gets the category of this network.
func (n *iNetwork) GetCategory() (NLMNetworkCategory, error) {
	res, err := n.idispatch.CallMethod("GetCategory")
	if err != nil {
		return -1, fmt.Errorf("failed to call GetCategory method: %v", err)
	}
	networkCategoryAny := res.Value()
	networkCategoryInt32, ok := networkCategoryAny.(int32)
	if !ok {
		return -1, fmt.Errorf("unexpected result type for GetCategory method: expected int32 but got %T", networkCategoryAny)
	}
	return NLMNetworkCategory(networkCategoryInt32), nil
}

// GetConnectivity gets the connectivity of this network.
func (n *iNetwork) GetConnectivity() (NLMConnectivity, error) {
	res, err := n.idispatch.CallMethod("GetConnectivity")
	if err != nil {
		return -1, fmt.Errorf("failed to call GetConnectivity method: %v", err)
	}
	networkConnectivityAny := res.Value()
	networkConnectivityInt32, ok := networkConnectivityAny.(int32)
	if !ok {
		return -1, fmt.Errorf("unexpected result type for GetConnectivity method: expected int32 but got %T", networkConnectivityAny)
	}
	return NLMConnectivity(networkConnectivityInt32), nil
}

// GetDescription gets the description/alias of this network.
func (n *iNetwork) GetDescription() (string, error) {
	res, err := n.idispatch.CallMethod("GetDescription")
	if err != nil {
		return "", fmt.Errorf("failed to call GetDescription method: %v", err)
	}
	return res.ToString(), nil
}

// GetDomainType gets the domain type of this network.
func (n *iNetwork) GetDomainType() (NLMDomainType, error) {
	res, err := n.idispatch.CallMethod("GetDomainType")
	if err != nil {
		return -1, fmt.Errorf("failed to call GetDomainType method: %v", err)
	}
	domainTypeAny := res.Value()
	domainTypeInt32, ok := domainTypeAny.(int32)
	if !ok {
		return -1, fmt.Errorf("unexpected result type for GetDomainType method: expected int32 but got %T", domainTypeAny)
	}
	return NLMDomainType(domainTypeInt32), nil
}

// GetName gets the name of this network.
func (n *iNetwork) GetName() (string, error) {
	res, err := n.idispatch.CallMethod("GetName")
	if err != nil {
		return "", fmt.Errorf("failed to call GetName method: %v", err)
	}
	return res.ToString(), nil
}

// GetNetworkConnections returns the network connections for this network.
func (n *iNetwork) GetNetworkConnections() (IEnumNetworkConnections, error) {
	res, err := n.idispatch.CallMethod("GetNetworkConnections")
	if err != nil {
		return nil, fmt.Errorf("failed to call GetNetworkConnections on INetwork object: %v", err)
	}
	idispatch := res.ToIDispatch()
	if idispatch == nil {
		return nil, fmt.Errorf("result of GetNetworkConnections is not an IDispatch, got type %T", res.Value())
	}
	defer idispatch.Release()
	networkConnections, err := NewNetworkConnections(idispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to get NetworkConnections from *ole.VARIANT: %v", err)
	}
	return networkConnections, nil
}

// SetCategory sets the category of this network.
func (n *iNetwork) SetCategory(category NLMNetworkCategory) (err error) {
	if _, err := n.idispatch.CallMethod("SetCategory", int32(category)); err != nil {
		return fmt.Errorf("failed to call SetCategory method with value %d: %v", int32(category), err)
	}
	return nil
}

// SetDescription sets the description/alias of this network.
func (n *iNetwork) SetDescription(descr string) error {
	if _, err := n.idispatch.CallMethod("SetDescription", descr); err != nil {
		return fmt.Errorf("failed to call SetDescription method with value %s: %v", descr, err)
	}
	return nil
}

// SetName sets the name of this network.
func (n *iNetwork) SetName(name string) error {
	if _, err := n.idispatch.CallMethod("SetName", name); err != nil {
		return fmt.Errorf("failed to call SetName method with value %s: %v", name, err)
	}
	return nil
}

// Release releases the INetwork object.
func (n *iNetwork) Release() {
	n.idispatch.Release()
}
