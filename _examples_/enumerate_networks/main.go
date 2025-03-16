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

	conns.ForEach(func(conn wnlm.INetworkConnection) bool {
		net, err := conn.GetNetwork()
		if err != nil {
			err = fmt.Errorf("failed to get network for conn")
			return false
		}
		defer net.Release()

		netName, err := net.GetName()
		if err != nil {
			err = fmt.Errorf("failed to get network name for network: %v", err)
			return false
		}
		log.Println(fmt.Sprintf("Name: \"%s\"", netName))

		netDesc, err := net.GetDescription()
		if err != nil {
			err = fmt.Errorf("failed to get network description for network %s: %v", netName, err)
			return false
		}
		log.Println(fmt.Sprintf("Description: \"%s\"", netDesc))

		netCategory, err := net.GetCategory()
		if err != nil {
			err = fmt.Errorf("failed to get network category on network %s: %v", netName, err)
			return false
		}
		log.Println(fmt.Sprintf("Category: \"%s\"", netCategory.String()))

		// make em public
		if netCategory == wnlm.NLMNetworkCategoryPrivate {
			net.SetCategory(wnlm.NLMNetworkCategoryPublic)
		}

		netDomainType, err := net.GetDomainType()
		if err != nil {
			err = fmt.Errorf("failed to get network domain type on network %s: %v", netName, err)
			return false
		}
		log.Println(fmt.Sprintf("Domain Type: \"%s\"", netDomainType.String()))

		netConnectivity, err := net.GetConnectivity()
		if err != nil {
			err = fmt.Errorf("failed to get network connectivity on network %s: %v", netName, err)
			return false
		}
		log.Println(fmt.Sprintf("Connectivity: \"%s\"", netConnectivity.String()))

		nconns, err := net.GetNetworkConnections()
		if err != nil {
			log.Printf("failed :%v", err)
			err = fmt.Errorf("failed to get network connections on network %s: %v", netName, err)
			return false
		}
		defer nconns.Release()
		log.Println(fmt.Sprintf("Network Connections: %d", nconns.Size()))

		guid, err := net.GetNetworkId()
		if err != nil {
			log.Printf("failed :%v", err)
			err = fmt.Errorf("failed to get network id on network %s: %v", netName, err)
			return false
		}
		log.Println(fmt.Sprintf("Network ID: %s", guid.String()))

		created, connected, err := net.GetTimeCreatedAndConnected()
		if err != nil {
			log.Printf("failed :%v", err)
			err = fmt.Errorf("failed to get network created/connected timestamps on network %s: %v", netName, err)
			return false
		}
		log.Println(fmt.Sprintf("Created At: %s", created.String()))
		log.Println(fmt.Sprintf("Connected At: %s", connected.String()))

		isConnectedToInternet, err := net.IsConnectedToInternet()
		if err != nil {
			log.Printf("failed :%v", err)
			err = fmt.Errorf("failed to get network isConnectedToInternet on network %s: %v", netName, err)
			return false
		}
		log.Println(fmt.Sprintf("Connected to internet: %t", isConnectedToInternet))

		isConnected, err := net.IsConnected()
		if err != nil {
			log.Printf("failed :%v", err)
			err = fmt.Errorf("failed to get network isConnected on network %s: %v", netName, err)
			return false
		}
		log.Println(fmt.Sprintf("Connected: %t", isConnected))

		log.Println()
		return true // continue
	})
	if err != nil {
		log.Fatalf("failed to iterate over network connections: %v", err)
	}
}
