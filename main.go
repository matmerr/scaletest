package main

import (
	"os"
)

const yamlDirectory = "./yaml/"

func init() {
	os.RemoveAll(yamlDirectory)
	//make yaml directory
	if err := os.MkdirAll(yamlDirectory, os.ModePerm); err != nil {
		panic(err)
	}
}

const (
	namespaces = 35

	serverDeploymentsPerNamespace = 5
	serverReplicasPerDeployment   = 150

	serverServicesPerNamespace = 5

	clientDeploymentsPerNamespace = 5
	clientReplicasPerDeployment   = 150
)

/*
1500 nodes
375 deployments
200 pods per deployment
40 namespaces
5 deployments per namespace
1000 pods per namespace
*/

func main() {

}
