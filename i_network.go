//go:build windows

package wnlm

import (
	"unsafe"

	"github.com/go-ole/go-ole"
)

// INetwork represents the Windows INetwork type as defined in
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/nn-netlistmgr-inetwork.
//
// The Windows Global Unique Identifier (GUID) for this interface is DCB00002-570F-4A9B-8D69-199FDBA5723B.
type INetwork interface {
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

// GetCategory get the category of a network.
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
