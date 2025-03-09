//go:build windows

package wnlm

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows"
)

// IEnumNetworkConnections represents an enumeration of the Windows INetworkConnections type as defined in
// https://learn.microsoft.com/en-us/windows/win32/api/netlistmgr/nn-netlistmgr-ienumnetworkconnections.
//
// The Windows Global Unique Identifier (GUID) for this interface is DCB00006-570F-4A9B-8D69-199FDBA5723B.
type IEnumNetworkConnections interface {
	FindInterfaceByGUID(guid *windows.GUID) (INetworkConnection, bool, error)
	Release()
}

// iEnumNetworkConnections is the default implementation of IEnumNetworkConnections.
type iEnumNetworkConnections struct {
	conns []INetworkConnection
}

// NewNetworkConnections returns an IEnumNetworkConnections object based on an IDispatch object.
func NewNetworkConnections(idispatch *ole.IDispatch) (IEnumNetworkConnections, error) {
	conns := []INetworkConnection{}
	err := oleutil.ForEach(idispatch, func(variant *ole.VARIANT) error {
		networkConnection, err := NewNetworkConnectionFromVariant(variant)
		if err != nil {
			return fmt.Errorf("failed to convert variant to NetworkConnection interface: %v", err)
		}
		conns = append(conns, networkConnection)
		return nil
	})
	if err != nil {
		// release any connections we had already fetched successfully
		for _, conn := range conns {
			conn.Release()
		}
		return nil, fmt.Errorf("failed to convert variant to list of network connections: %v", err)
	}
	return &iEnumNetworkConnections{conns: conns}, nil
}

// FindInterfaceByGUID returns the INetworkConnection object for a network interface by GUID.
func (nc *iEnumNetworkConnections) FindInterfaceByGUID(guid *windows.GUID) (INetworkConnection, bool, error) {
	for _, networkConnection := range nc.conns {
		adapterGUID, err := networkConnection.GetAdapterId()
		if err != nil {
			return nil, false, fmt.Errorf("failed to get adapter id for network connection: %v", err)
		}
		if adapterGUID.String() != guid.String() {
			continue
		}
		return networkConnection, true, nil
	}
	return nil, false, nil
}

// Release releases the IEnumNetworkConnections object.
func (nc *iEnumNetworkConnections) Release() {
	for _, conn := range nc.conns {
		conn.Release()
	}
}
