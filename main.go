package main

import (
	"fmt"
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
	serverReplicasPerDeployment   = 5

	serverServicesPerNamespace = 5

	clientDeploymentsPerNamespace = 5
	clientReplicasPerDeployment   = 2
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
	CreateNamespaces(yamlDirectory, namespaces)

	for nsNum := 0; nsNum < namespaces; nsNum++ {
		targetDirectory := fmt.Sprintf("%s/%d", yamlDirectory, nsNum)

		// create all deploymens in the namespace
		for deployNum := 0; deployNum < serverDeploymentsPerNamespace; deployNum++ {
			// create server deployment
			server := FortioServerDeployment{
				Name:                "fortio-server-" + fmt.Sprint(deployNum),
				Namespace:           "fortio-" + fmt.Sprint(nsNum),
				Replicas:            serverReplicasPerDeployment,
				ServiceBackendLabel: "fortio-service-" + fmt.Sprint(nsNum),
				AppLabel:            "fortio-server-" + fmt.Sprint(deployNum),
				NodeSelector:        "scenario: podcount",
			}
			CreateServerDeployments(fmt.Sprintf("%s/1-%d-server.yaml", targetDirectory, deployNum), &server)
		}

		// create the
		for svcNum := 0; svcNum < serverServicesPerNamespace; svcNum++ {
			service := FortioService{
				Name:                "fortio-service-" + fmt.Sprint(svcNum),
				Namespace:           "fortio-" + fmt.Sprint(nsNum),
				TargetPort:          "8080",
				ServiceBackendLabel: "fortio-service-" + fmt.Sprint(nsNum),
			}
			CreateServiceDeployments(fmt.Sprintf("%s/2-%d-service.yaml", targetDirectory, svcNum), &service)
		}

		// create all Client deployments
		for clientNum := 0; clientNum < serverServicesPerNamespace; clientNum++ {

			svcNum := clientNum % serverServicesPerNamespace

			client := FortioClientDeployment{
				Name:         "fortio-client-" + fmt.Sprint(clientNum),
				Namespace:    "fortio-" + fmt.Sprint(nsNum),
				Replicas:     clientReplicasPerDeployment,
				RequestURL:   "fortio-service-" + fmt.Sprint(svcNum),
				RequestPort:  "8080",
				AppLabel:     "fortio-client-" + fmt.Sprint(clientNum),
				QPS:          "2500",
				NodeSelector: "scenario: podcount",
			}
			CreateClientDeployments(fmt.Sprintf("%s/3-%d-client.yaml", targetDirectory, clientNum), &client)
		}
	}
}
