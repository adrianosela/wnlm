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
			err = fmt.Errorf("failed to get network for conn: %v", err2)
			return false
		}
		defer net.Release()

		netName, err2 := net.GetName()
		if err2 != nil {
			err = fmt.Errorf("failed to get network name for network: %v", err2)
			return false
		}
		fmt.Printf("Name: \"%s\"\n", netName)

		netDesc, err2 := net.GetDescription()
		if err2 != nil {
			err = fmt.Errorf("failed to get network description for network %s: %v", netName, err2)
			return false
		}
		fmt.Printf("Description: \"%s\"\n", netDesc)

		netCategory, err2 := net.GetCategory()
		if err2 != nil {
			err = fmt.Errorf("failed to get network category on network %s: %v", netName, err2)
			return false
		}
		fmt.Printf("Category: %s\n", netCategory.String())

		netDomainType, err2 := net.GetDomainType()
		if err2 != nil {
			err = fmt.Errorf("failed to get network domain type on network %s: %v", netName, err2)
			return false
		}
		fmt.Printf("Domain Type: %s\n", netDomainType.String())

		netConnectivity, err2 := net.GetConnectivity()
		if err2 != nil {
			err = fmt.Errorf("failed to get network connectivity on network %s: %v", netName, err2)
			return false
		}
		fmt.Printf("Connectivity: %s\n", netConnectivity.String())

		nconns, err2 := net.GetNetworkConnections()
		if err2 != nil {
			err = fmt.Errorf("failed to get network connections on network %s: %v", netName, err2)
			return false
		}
		defer nconns.Release()
		fmt.Printf("Network Connections: %d\n", nconns.Size())

		guid, err2 := net.GetNetworkId()
		if err2 != nil {
			err = fmt.Errorf("failed to get network id on network %s: %v", netName, err2)
			return false
		}
		fmt.Printf("Network ID: %s\n", guid.String())

		created, connected, err2 := net.GetTimeCreatedAndConnected()
		if err2 != nil {
			err = fmt.Errorf("failed to get network created/connected timestamps on network %s: %v", netName, err2)
			return false
		}
		fmt.Printf("Created At: %s\n", created.String())
		fmt.Printf("Connected At: %s\n", connected.String())

		isConnectedToInternet, err2 := net.IsConnectedToInternet()
		if err2 != nil {
			err = fmt.Errorf("failed to get network isConnectedToInternet on network %s: %v", netName, err2)
			return false
		}
		fmt.Printf("Connected to Internet: %t\n", isConnectedToInternet)

		isConnected, err2 := net.IsConnected()
		if err2 != nil {
			err = fmt.Errorf("failed to get network isConnected on network %s: %v", netName, err2)
			return false
		}
		fmt.Printf("Connected: %t\n", isConnected)

		return true // continue
	})
	if err != nil {
		log.Fatalf("failed while iterating over network connections: %v", err)
	}
}
