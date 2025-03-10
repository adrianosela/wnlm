//go:build windows

package wnlm

import (
	"fmt"
	"unsafe"

	"github.com/go-ole/go-ole"
)

// INetwork represents the Windows INetwork type as defined in
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/nn-netlistmgr-inetwork.
//
// The Windows Global Unique Identifier (GUID) for this interface is DCB00002-570F-4A9B-8D69-199FDBA5723B.
type INetwork interface {
	GetName() (string, error)
	GetDescription() (string, error)
	GetDomainType() (NLMDomainType, error)
	SetDomainType(NLMDomainType) error
	GetCategory() (NLMNetworkCategory, error)
	SetCategory(NLMNetworkCategory) error
	Release()
}

// iNetwork is the default implementation of INetwork.
type iNetwork struct {
	idispatch *ole.IDispatch
}

// iNetworkVTable represents the INetwork interface's VTable.
type iNetworkVTable struct {
	ole.IDispatchVtbl
	GetName                    uintptr
	SetName                    uintptr
	GetDescription             uintptr
	SetDescription             uintptr
	GetNetworkId               uintptr
	GetDomainType              uintptr
	GetNetworkConnections      uintptr
	GetTimeCreatedAndConnected uintptr
	Get_IsConnectedToInternet  uintptr
	Get_IsConnected            uintptr
	GetConnectivity            uintptr
	GetCategory                uintptr
	SetCategory                uintptr
}

// GetName gets the name of a network.
func (n *iNetwork) GetName() (string, error) {
	resp, err := n.idispatch.CallMethod("GetName")
	if err != nil {
		return "", fmt.Errorf("failed to call GetName method: %v", err)
	}
	return resp.ToString(), nil
}

// GetDescription gets the description of a network.
func (n *iNetwork) GetDescription() (string, error) {
	resp, err := n.idispatch.CallMethod("GetDescription")
	if err != nil {
		return "", fmt.Errorf("failed to call GetDescription method: %v", err)
	}
	return resp.ToString(), nil
}

// GetDomainType gets the domain type of a network.
func (n *iNetwork) GetDomainType() (NLMDomainType, error) {
	var domainType byte
	hr, _, _ := globalSyscaller.SyscallN(
		(*iNetworkVTable)(unsafe.Pointer(n.idispatch.RawVTable)).GetDomainType,
		uintptr(unsafe.Pointer(n.idispatch)),
		uintptr(unsafe.Pointer(&domainType)),
	)
	if hr < 0 {
		return 0, ole.NewError(hr)
	}
	return NLMDomainType(domainType), nil
}

// SetDomainType sets the domain type on a network.
func (n *iNetwork) SetDomainType(domainType NLMDomainType) (err error) {
	hr, _, _ := globalSyscaller.SyscallN(
		(*iNetworkVTable)(unsafe.Pointer(n.idispatch.RawVTable)).SetCategory,
		uintptr(unsafe.Pointer(n.idispatch)),
		uintptr(domainType))
	if hr < 0 {
		return ole.NewError(hr)
	}
	return nil
}

// GetCategory gets the category of a network.
func (n *iNetwork) GetCategory() (NLMNetworkCategory, error) {
	var category byte
	hr, _, _ := globalSyscaller.SyscallN(
		(*iNetworkVTable)(unsafe.Pointer(n.idispatch.RawVTable)).GetCategory,
		uintptr(unsafe.Pointer(n.idispatch)),
		uintptr(unsafe.Pointer(&category)),
	)
	if hr < 0 {
		return 0, ole.NewError(hr)
	}
	return NLMNetworkCategory(category), nil
}

// SetCategory sets the category on a network.
func (n *iNetwork) SetCategory(category NLMNetworkCategory) (err error) {
	hr, _, _ := globalSyscaller.SyscallN(
		(*iNetworkVTable)(unsafe.Pointer(n.idispatch.RawVTable)).SetCategory,
		uintptr(unsafe.Pointer(n.idispatch)),
		uintptr(category))
	if hr < 0 {
		return ole.NewError(hr)
	}
	return nil
}

// Release releases the INetwork object.
func (n *iNetwork) Release() {
	n.idispatch.Release()
}
