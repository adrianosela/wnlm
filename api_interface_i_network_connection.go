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
	GetNetwork() (INetwork, error)
	IsConnectedToInternet() (bool, error)
	IsConnected() (bool, error)
	GetConnectivity() (NLMConnectivity, error)
	GetConnectionId() (*ole.GUID, error)
	GetAdapterId() (*ole.GUID, error)
	GetDomainType() (NLMDomainType, error)

	Release()
}

// iNetworkConnections is the default implementation of INetworkConnection.
type iNetworkConnection struct {
	idispatch *ole.IDispatch
}

// iNetworkConnectionVTable represents the INetworkConnection interface's VTable.
type iNetworkConnectionVTable struct {
	ole.IDispatchVtbl
	GetNetwork            uintptr // id = 1, method
	IsConnectedToInternet uintptr // id = 2, property
	IsConnected           uintptr // id = 3, property
	GetConnectivity       uintptr // id = 4, method
	GetConnectionId       uintptr // id = 5, method
	GetAdapterId          uintptr // id = 6, method
	GetDomainType         uintptr // id = 7, method
}

// vtable returns the INetworkConnection's VTable.
func (n *iNetworkConnection) vtable() *iNetworkConnectionVTable {
	return (*iNetworkConnectionVTable)(unsafe.Pointer(n.idispatch.RawVTable))
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

// IsConnectedToInternet returns whether the network connection is connected to the Internet.
func (n *iNetworkConnection) IsConnectedToInternet() (bool, error) {
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

// IsConnected returns whether the network connection is connected.
func (n *iNetworkConnection) IsConnected() (bool, error) {
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

// GetNetwork returns the INetwork for a network connection.
func (nc *iNetworkConnection) GetNetwork() (INetwork, error) {
	var idispatch *ole.IDispatch
	hr, _, _ := syscall.SyscallN(
		nc.vtable().GetNetwork,
		uintptr(unsafe.Pointer(nc.idispatch)),
		uintptr(unsafe.Pointer(&idispatch)),
	)
	if hr < 0 {
		return nil, ole.NewError(hr)
	}
	return INetworkFromIDispatch(idispatch), nil
}

// GetConnectivity gets the connectivity of this network connection.
func (n *iNetworkConnection) GetConnectivity() (NLMConnectivity, error) {
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

// GetConnectionId returns the connection GUID for a network connection.
func (nc *iNetworkConnection) GetConnectionId() (*ole.GUID, error) {
	var guid ole.GUID
	hr, _, _ := syscall.SyscallN(
		nc.vtable().GetConnectionId,
		uintptr(unsafe.Pointer(nc.idispatch)),
		uintptr(unsafe.Pointer(&guid)),
	)
	if hr < 0 {
		return nil, fmt.Errorf("failed to get connection id for network connection: %v", ole.NewError(hr))
	}
	return &guid, nil
}

// GetAdapterId returns the adapter GUID for a network connection.
func (nc *iNetworkConnection) GetAdapterId() (*ole.GUID, error) {
	var guid ole.GUID
	hr, _, _ := syscall.SyscallN(
		nc.vtable().GetAdapterId,
		uintptr(unsafe.Pointer(nc.idispatch)),
		uintptr(unsafe.Pointer(&guid)),
	)
	if hr < 0 {
		return nil, fmt.Errorf("failed to get adapter id for network connection: %v", ole.NewError(hr))
	}
	return &guid, nil
}

// GetDomainType gets the domain type of this network connection.
func (n *iNetworkConnection) GetDomainType() (NLMDomainType, error) {
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

// Release releases the INetworkConnection object.
func (nc *iNetworkConnection) Release() {
	nc.idispatch.Release()
}
