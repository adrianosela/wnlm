//go:build windows

package wnlm

import (
	"fmt"

	"github.com/go-ole/go-ole"
)

// INetworkListManager represents the Windows INetworkListManager type as defined in
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/nn-netlistmgr-inetworklistmanager.
//
// The Windows Global Unique Identifier (GUID) for this interface is DCB00000-570F-4A9B-8D69-199FDBA5723B.
type INetworkListManager interface {
	GetNetworkConnections() (IEnumNetworkConnections, error)

	Release()
}

// iNetworkListManager is the default implementation of INetworkListManager.
type iNetworkListManager struct {
	idispatch *ole.IDispatch
}

// NewNetworkListManager initializes a new Windows NetworkListManager API client.
func NewNetworkListManager() (INetworkListManager, error) {
	// NOTE(@adrianosela): DCB00C01-570F-4A9B-8D69-199FDBA5723B is the
	// well-known Windows Global ID for the NetworkListManager class ID.
	if err := globalOleConn.Create("{DCB00C01-570F-4A9B-8D69-199FDBA5723B}"); err != nil {
		return nil, fmt.Errorf("failed to create NetworkListManager object by program GUID: %v", err)
	}
	defer globalOleConn.Release()
	dispatch, err := globalOleConn.Dispatch()
	if err != nil {
		return nil, fmt.Errorf("failed to get dispatch object for NetworkListManager: %v", err)
	}
	return &iNetworkListManager{idispatch: dispatch.Object}, nil
}

// GetNetworkConnections returns all network connections for the system by using the API call described in
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/nf-netlistmgr-inetwork-getnetworkconnections.
func (nlm *iNetworkListManager) GetNetworkConnections() (IEnumNetworkConnections, error) {
	networkConnectionsVariant, err := nlm.idispatch.CallMethod("GetNetworkConnections")
	if err != nil {
		return nil, fmt.Errorf("failed to call GetNetworkConnections on NetworkListManager object: %v", err)
	}
	idispatch := networkConnectionsVariant.ToIDispatch()
	if idispatch == nil {
		return nil, fmt.Errorf("failed to convert variant to IDispatch")
	}
	defer idispatch.Release()

	networkConnections, err := NewNetworkConnections(idispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to get NetworkConnections from *ole.VARIANT: %v", err)
	}
	return networkConnections, nil
}

// Release releases the NetworkListManager object.
func (nlm *iNetworkListManager) Release() {
	nlm.idispatch.Release()
}
