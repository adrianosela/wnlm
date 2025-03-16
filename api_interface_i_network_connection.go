//go:build windows

package wnlm

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

// INetworkConnection represents the Windows INetworkConnection type as defined in
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/nn-netlistmgr-inetworkconnection.
//
// The Windows Global Unique Identifier (GUID) for this interface is DCB00005-570F-4A9B-8D69-199FDBA5723B.
type INetworkConnection interface {
	GetAdapterId() (*ole.GUID, error)
	GetNetwork() (INetwork, error)

	Release()
}

// iNetworkConnections is the default implementation of INetworkConnection.
type iNetworkConnection struct {
	idispatch *ole.IDispatch
}

// iNetworkConnectionVTable represents the INetworkConnection interface's VTable.
type iNetworkConnectionVTable struct {
	ole.IDispatchVtbl
	GetNetwork                uintptr
	Get_IsConnectedToInternet uintptr
	Get_IsConnected           uintptr
	GetConnectivity           uintptr
	GetConnectionId           uintptr
	GetAdapterId              uintptr
	GetDomainType             uintptr
}

// NewNetworkConnectionFromVariant returns the INetworkConnection object for a given ole.VARIANT.
func NewNetworkConnectionFromVariant(variant *ole.VARIANT) (INetworkConnection, error) {
	iUnknown := variant.ToIUnknown()
	if iUnknown == nil {
		return nil, fmt.Errorf("expected variant to be of VT type %d, but got %d", ole.VT_UNKNOWN, variant.VT)
	}

	// NOTE(@adrianosela): {DCB00005-570F-4A9B-8D69-199FDBA5723B} is the
	// well-known Windows Global ID for the INetworkConnection interface.
	interfaceGUID := ole.NewGUID("{DCB00005-570F-4A9B-8D69-199FDBA5723B}")

	idispatch, err := iUnknown.QueryInterface(interfaceGUID)
	if err != nil {
		return nil, fmt.Errorf("failed to use unknown interface as interface with GUID %s: %v", interfaceGUID.String(), err)
	}
	return &iNetworkConnection{idispatch: idispatch}, nil
}

// GetAdapterId returns the adapter GUID for a network connection.
func (nc *iNetworkConnection) GetAdapterId() (*ole.GUID, error) {
	var guid ole.GUID
	hr, _, _ := syscall.SyscallN(
		(*iNetworkConnectionVTable)(unsafe.Pointer(nc.idispatch.RawVTable)).GetAdapterId,
		uintptr(unsafe.Pointer(nc.idispatch)),
		uintptr(unsafe.Pointer(&guid)),
	)
	if hr < 0 {
		return nil, fmt.Errorf("failed to get adapter id for network connection: %v", ole.NewError(hr))
	}
	return &guid, nil
}

// GetAdapterId returns the INetwork for a network connection.
func (nc *iNetworkConnection) GetNetwork() (INetwork, error) {
	var idispatch *ole.IDispatch
	hr, _, _ := syscall.SyscallN(
		(*iNetworkConnectionVTable)(unsafe.Pointer(nc.idispatch.RawVTable)).GetNetwork,
		uintptr(unsafe.Pointer(nc.idispatch)),
		uintptr(unsafe.Pointer(&idispatch)),
	)
	if hr < 0 {
		return nil, ole.NewError(hr)
	}
	return INetworkFromIDispatch(idispatch), nil
}

// Release releases the INetworkConnection object.
func (nc *iNetworkConnection) Release() {
	nc.idispatch.Release()
}
