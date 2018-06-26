package basic

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/perlin-network/noise/examples/basic/messages"
)

func ExampleSetupClusters() {
	// parse to flags to silence the glog library
	flag.Parse()

	host := "localhost"
	cluster1StartPort := 5000
	cluster1NumPorts := 3
	nodes := []*ClusterNode{}

	for i := 0; i < cluster1NumPorts; i++ {
		node := &ClusterNode{}
		node.Host = host
		node.Port = cluster1StartPort + i

		nodes = append(nodes, node)
	}

	if err := SetupCluster(nodes); err != nil {
		fmt.Print(err)
	}

	for i, node := range nodes {
		if node.Net == nil {
			fmt.Printf("Expected %d nodes, but node %d is missing a network", len(nodes), i)
		}
	}

	// check if you can send a message from node 1 and will it be received only in node 2,3
	{
		testMessage := "message from node 0"
		nodes[0].Net.Broadcast(&messages.ClusterTestMessage{Message: testMessage})

		// HACK: TODO: replace sleep with something else
		time.Sleep(1 * time.Second)

		if result := nodes[0].PopMessage(); result != nil {
			glog.Errorf("Expected nothing in node 0, got %v", result)
		}
		for i := 1; i < len(nodes); i++ {
			if result := nodes[i].PopMessage(); result == nil {
				fmt.Printf("Expected a message in node %d but it was blank", i)
			} else {
				if result.Message != testMessage {
					fmt.Printf("Expected message %s in node %d but got %v", testMessage, i, result)
				}
			}
		}
	}

	fmt.Printf("Success")
	// Output: Success
}
