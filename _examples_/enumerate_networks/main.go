//go:build windows

package main

import (
	"fmt"
	"log"

	"github.com/adrianosela/wnlm"
)

func main() {
	wnlm.Initialize()
	defer wnlm.Uninitialize()

	nlm, err := wnlm.NewNetworkListManager()
	if err != nil {
		log.Fatalf("failed to initialize network list manager: %v", err)
	}
	defer nlm.Release()

	conns, err := nlm.GetNetworkConnections()
	if err != nil {
		log.Fatalf("failed to list network connections: %v", err)
	}
	defer conns.Release()

	conns.ForEach(func(i int, conn wnlm.INetworkConnection) bool {
		if i > 0 {
			fmt.Println()
		}

		net, err2 := conn.GetNetwork()
		if err2 != nil {
			err = fmt.Errorf("failed to get INetwork for INetworkConnection: %v", err2)
			return false
		}
		defer net.Release()

		adapterID, err2 := conn.GetAdapterId()
		if err2 != nil {
			err = fmt.Errorf("failed to get adapter ID for INetworkConnection: %v", err2)
			return false
		}
		fmt.Printf("INetworkConnection Adapter ID: \"%s\"\n", adapterID.String())

		connID, err2 := conn.GetConnectionId()
		if err2 != nil {
			err = fmt.Errorf("failed to get connection ID for INetworkConnection: %v", err2)
			return false
		}
		fmt.Printf("INetworkConnection Connection ID: \"%s\"\n", connID.String())

		isConnectedToInternetConn, err2 := conn.IsConnectedToInternet()
		if err2 != nil {
			err = fmt.Errorf("failed to get isConnectedToInternet on INetworkConnection %s", err2)
			return false
		}
		fmt.Printf("INetworkConnection Connected to Internet: %t\n", isConnectedToInternetConn)

		isConnectedConn, err2 := conn.IsConnected()
		if err2 != nil {
			err = fmt.Errorf("failed to get isConnected on INetworkConnection: %v", err2)
			return false
		}
		fmt.Printf("INetworkConnection Connected: %t\n", isConnectedConn)

		netName, err2 := net.GetName()
		if err2 != nil {
			err = fmt.Errorf("failed to get name for INetwork: %v", err2)
			return false
		}
		fmt.Printf("INetwork Name: \"%s\"\n", netName)

		netDesc, err2 := net.GetDescription()
		if err2 != nil {
			err = fmt.Errorf("failed to get description for INetwork %s: %v", netName, err2)
			return false
		}
		fmt.Printf("INetwork Description: \"%s\"\n", netDesc)

		netCategory, err2 := net.GetCategory()
		if err2 != nil {
			err = fmt.Errorf("failed to get category for INetwork %s: %v", netName, err2)
			return false
		}
		fmt.Printf("INetwork Category: %s\n", netCategory.String())

		netDomainType, err2 := net.GetDomainType()
		if err2 != nil {
			err = fmt.Errorf("failed to get domain type for INetwork %s: %v", netName, err2)
			return false
		}
		fmt.Printf("INetwork Domain Type: %s\n", netDomainType.String())

		netConnectivity, err2 := net.GetConnectivity()
		if err2 != nil {
			err = fmt.Errorf("failed to get connectivity for INetwork %s: %v", netName, err2)
			return false
		}
		fmt.Printf("INetwork Connectivity: %s\n", netConnectivity.String())

		nconns, err2 := net.GetNetworkConnections()
		if err2 != nil {
			err = fmt.Errorf("failed to get IEnumNetworkConnections for INetwork %s: %v", netName, err2)
			return false
		}
		defer nconns.Release()
		fmt.Printf("INetwork (# of) Network Connections: %d\n", nconns.Size())

		guid, err2 := net.GetNetworkId()
		if err2 != nil {
			err = fmt.Errorf("failed to get id for INetwork %s: %v", netName, err2)
			return false
		}
		fmt.Printf("INetwork Network ID: %s\n", guid.String())

		created, connected, err2 := net.GetTimeCreatedAndConnected()
		if err2 != nil {
			err = fmt.Errorf("failed to get created/connected timestamps for INetwork %s: %v", netName, err2)
			return false
		}
		fmt.Printf("INetwork Created At: %s\n", created.String())
		fmt.Printf("INetwork Connected At: %s\n", connected.String())

		isConnectedToInternet, err2 := net.IsConnectedToInternet()
		if err2 != nil {
			err = fmt.Errorf("failed to get isConnectedToInternet on INetwork %s: %v", netName, err2)
			return false
		}
		fmt.Printf("INetwork Connected to Internet: %t\n", isConnectedToInternet)

		isConnected, err2 := net.IsConnected()
		if err2 != nil {
			err = fmt.Errorf("failed to get isConnected on INetwork %s: %v", netName, err2)
			return false
		}
		fmt.Printf("INetwork Connected: %t\n", isConnected)

		return true // continue
	})
	if err != nil {
		log.Fatalf("failed while iterating over IEnumNetworkConnection: %v", err)
	}
}
