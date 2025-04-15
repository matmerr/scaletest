package steps

import (
	"context"
	"fmt"

	"github.com/matmerr/scaletest/yaml"
	"github.com/matmerr/scaletest/yaml/fortio"
)

type GenerateYamlsStep struct {
	Directory                     string
	Namespaces                    int
	ServerDeploymentsPerNamespace int
	ServerReplicasPerDeployment   int
	ServerServicesPerNamespace    int
	ClientDeploymentsPerNamespace int
	ClientReplicasPerDeployment   int
}

func (g *GenerateYamlsStep) Do(ctx context.Context) error {
	for nsNum := 0; nsNum < g.Namespaces; nsNum++ {
		targetDirectory := fmt.Sprintf("%s/longrunning/%d", g.Directory, nsNum)

		namespace := fortio.Namespace{
			Name: "fortio-" + fmt.Sprint(nsNum),
		}
		err := yaml.CreateYamlFile(fmt.Sprintf("%s/0-%d-ns.yaml", targetDirectory, nsNum), &namespace)
		if err != nil {
			return fmt.Errorf("failed to create namespace yaml file: %w", err)
		}

		// create all deploymens in the namespace
		for deployNum := 0; deployNum < g.ServerDeploymentsPerNamespace; deployNum++ {
			// create server deployment
			server := fortio.FortioServerDeployment{
				Name:                "fortio-server-" + fmt.Sprint(deployNum),
				Namespace:           "fortio-" + fmt.Sprint(nsNum),
				Replicas:            g.ServerReplicasPerDeployment,
				ServiceBackendLabel: "fortio-service-" + fmt.Sprint(nsNum),
				AppLabel:            "fortio-server-" + fmt.Sprint(deployNum),
				NodeSelector:        "scenario: podcount",
			}

			err := yaml.CreateYamlFile(fmt.Sprintf("%s/1-%d-server.yaml", targetDirectory, deployNum), &server)
			if err != nil {
				return fmt.Errorf("failed to create server yaml file: %w", err)
			}
		}

		// create the services in the namespace
		for svcNum := 0; svcNum < g.ServerServicesPerNamespace; svcNum++ {
			service := fortio.FortioService{
				Name:                "fortio-service-" + fmt.Sprint(svcNum),
				Namespace:           "fortio-" + fmt.Sprint(nsNum),
				TargetPort:          "8080",
				ServiceBackendLabel: "fortio-service-" + fmt.Sprint(nsNum),
			}

			err := yaml.CreateYamlFile(fmt.Sprintf("%s/2-%d-service.yaml", targetDirectory, svcNum), &service)
			if err != nil {
				return fmt.Errorf("failed to create service yaml file: %w", err)
			}
		}

		// create all client deployments
		for clientNum := 0; clientNum < g.ClientDeploymentsPerNamespace; clientNum++ {

			svcNum := clientNum % g.ServerServicesPerNamespace

			client := fortio.FortioClientDeployment{
				Name:         "fortio-client-" + fmt.Sprint(clientNum),
				Namespace:    "fortio-" + fmt.Sprint(nsNum),
				Replicas:     g.ClientReplicasPerDeployment,
				RequestURL:   "fortio-service-" + fmt.Sprint(svcNum),
				RequestPort:  "8080",
				AppLabel:     "fortio-client-" + fmt.Sprint(clientNum),
				QPS:          "2500",
				NodeSelector: "scenario: podcount",
			}
			err := yaml.CreateYamlFile(fmt.Sprintf("%s/3-%d-client.yaml", targetDirectory, clientNum), &client)
			if err != nil {
				return fmt.Errorf("failed to create client yaml file: %w", err)
			}
		}
	}

	return nil
}
