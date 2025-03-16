//go:build windows

package wnlm

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/adrianosela/wnlm/pkg/wintime"
	"github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

// INetwork represents the Windows INetwork type as defined in
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/nn-netlistmgr-inetwork.
//
// The Windows Global Unique Identifier (GUID) for this interface is DCB00002-570F-4A9B-8D69-199FDBA5723B.
type INetwork interface {
	GetName() (string, error)
	SetName(string) error
	GetDescription() (string, error)
	SetDescription(string) error
	GetNetworkId() (*windows.GUID, error)
	GetDomainType() (NLMDomainType, error)
	GetNetworkConnections() (IEnumNetworkConnections, error)
	GetTimeCreatedAndConnected() (time.Time, time.Time, error)
	IsConnectedToInternet() (bool, error)
	IsConnected() (bool, error)
	GetConnectivity() (NLMConnectivity, error)
	GetCategory() (NLMNetworkCategory, error)
	SetCategory(NLMNetworkCategory) error

	Release()
}

// iNetwork is the default implementation of INetwork.
type iNetwork struct {
	idispatch *ole.IDispatch
}

type iNetworkVtbl struct {
	ole.IDispatchVtbl
	GetName                    uintptr // id = 1, method
	SetName                    uintptr // id = 2, method
	GetDescription             uintptr // id = 3, method
	SetDescription             uintptr // id = 4, method
	GetNetworkId               uintptr // id = 5, method
	GetDomainType              uintptr // id = 6, method
	GetNetworkConnections      uintptr // id = 7, method
	GetTimeCreatedAndConnected uintptr // id = 8, method
	IsConnectedToInternet      uintptr // id = 9, property
	IsConnected                uintptr // id = 10, property
	GetConnectivity            uintptr // id = 11, method
	GetCategory                uintptr // id = 12, method
	SetCategory                uintptr // id = 13, method
}

// vtable returns the INetwork's VTable.
func (n *iNetwork) vtable() *iNetworkVtbl {
	return (*iNetworkVtbl)(unsafe.Pointer(n.idispatch.RawVTable))
}

// INetworkFromIDispatch returns an INetwork based on its IDispatch.
func INetworkFromIDispatch(idispatch *ole.IDispatch) INetwork {
	return &iNetwork{idispatch: idispatch}
}

// GetName gets the name of this network.
func (n *iNetwork) GetName() (string, error) {
	res, err := n.idispatch.CallMethod("GetName")
	if err != nil {
		return "", fmt.Errorf("failed to call GetName method: %v", err)
	}
	return res.ToString(), nil
}

// SetName sets the name of this network.
func (n *iNetwork) SetName(name string) error {
	if _, err := n.idispatch.CallMethod("SetName", name); err != nil {
		return fmt.Errorf("failed to call SetName method with value %s: %v", name, err)
	}
	return nil
}

// GetDescription gets the description/alias of this network.
func (n *iNetwork) GetDescription() (string, error) {
	res, err := n.idispatch.CallMethod("GetDescription")
	if err != nil {
		return "", fmt.Errorf("failed to call GetDescription method: %v", err)
	}
	return res.ToString(), nil
}

// SetDescription sets the description/alias of this network.
func (n *iNetwork) SetDescription(descr string) error {
	if _, err := n.idispatch.CallMethod("SetDescription", descr); err != nil {
		return fmt.Errorf("failed to call SetDescription method with value %s: %v", descr, err)
	}
	return nil
}

// GetNetworkId returns the GUID od this network.
func (n *iNetwork) GetNetworkId() (*windows.GUID, error) {
	guid := windows.GUID{}
	hr, _, _ := syscall.SyscallN(
		n.vtable().GetNetworkId,
		uintptr(unsafe.Pointer(n.idispatch)),
		uintptr(unsafe.Pointer(&guid)),
	)
	if hr != 0 {
		return nil, ole.NewError(hr)
	}
	return &guid, nil
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

// GetTimeCreatedAndConnected gets the timestamps of this network being created and connected.
func (n *iNetwork) GetTimeCreatedAndConnected() (time.Time, time.Time, error) {
	var createdLow, createdHigh, connectedLow, connectedHigh int64
	hr, _, _ := syscall.SyscallN(
		n.vtable().GetTimeCreatedAndConnected,
		uintptr(unsafe.Pointer(n.idispatch)),
		uintptr(unsafe.Pointer(&createdLow)),
		uintptr(unsafe.Pointer(&createdHigh)),
		uintptr(unsafe.Pointer(&connectedLow)),
		uintptr(unsafe.Pointer(&connectedHigh)),
	)
	if hr != 0 {
		return time.Time{}, time.Time{}, ole.NewError(hr)
	}
	created := wintime.ToTime(createdLow, createdHigh)
	connected := wintime.ToTime(connectedLow, connectedHigh)
	return created, connected, nil
}

// IsConnectedToInternet returns whether the network is connected to the Internet.
func (n *iNetwork) IsConnectedToInternet() (bool, error) {
	res, err := n.idispatch.GetProperty("IsConnectedToInternet")
	if err != nil {
		return false, fmt.Errorf("failed to get IsConnectedToInternet property: %v", err)
	}
	isConnectedToInternetAny := res.Value()
	isConnectedToInternetBool, ok := isConnectedToInternetAny.(bool)
	if !ok {
		return false, fmt.Errorf("unexpected result type for IsConnectedToInternet property: expected bool but got %T", isConnectedToInternetAny)
	}
	return isConnectedToInternetBool, nil
}

// IsConnected returns whether the network is connected.
func (n *iNetwork) IsConnected() (bool, error) {
	res, err := n.idispatch.GetProperty("IsConnected")
	if err != nil {
		return false, fmt.Errorf("failed to get IsConnected property: %v", err)
	}
	isConnectedAny := res.Value()
	isConnectedBool, ok := isConnectedAny.(bool)
	if !ok {
		return false, fmt.Errorf("unexpected result type for IsConnected property: expected bool but got %T", isConnectedAny)
	}
	return isConnectedBool, nil
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

// SetCategory sets the category of this network.
func (n *iNetwork) SetCategory(category NLMNetworkCategory) (err error) {
	if _, err := n.idispatch.CallMethod("SetCategory", int32(category)); err != nil {
		return fmt.Errorf("failed to call SetCategory method with value %d: %v", int32(category), err)
	}
	return nil
}

// Release releases the INetwork object.
func (n *iNetwork) Release() {
	n.idispatch.Release()
}
