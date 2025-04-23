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
	ClientQPS                     int
}

func (g *GenerateYamlsStep) Do(ctx context.Context) error {
	for nsNum := 0; nsNum < g.Namespaces; nsNum++ {

		namespace := fortio.Namespace{
			Name: "ns" + fmt.Sprint(nsNum),
		}

		nsDirectory := fmt.Sprintf("%s/longrunning/00-namespaces", g.Directory)
		err := yaml.CreateYamlFile(fmt.Sprintf("%s/0-%d-ns.yaml", nsDirectory, nsNum), &namespace)
		if err != nil {
			return fmt.Errorf("failed to create namespace yaml file: %w", err)
		}

		targetDirectory := fmt.Sprintf("%s/longrunning/ns-%d", g.Directory, nsNum)

		// create all deploymens in the namespace
		for deployNum := 0; deployNum < g.ServerDeploymentsPerNamespace; deployNum++ {
			server := fortio.FortioServerDeployment{
				Name:                fmt.Sprintf("ns%d-server-%d", nsNum, deployNum),
				Namespace:           "ns" + fmt.Sprint(nsNum),
				Replicas:            g.ServerReplicasPerDeployment,
				ServiceBackendLabel: "fortio-service-" + fmt.Sprint(deployNum),
				AppLabel:            "fortio-server-" + fmt.Sprint(deployNum),
				NodeSelector:        "scenario: highcount",
			}

			err := yaml.CreateYamlFile(fmt.Sprintf("%s/1-%d-server.yaml", targetDirectory, deployNum), &server)
			if err != nil {
				return fmt.Errorf("failed to create server yaml file: %w", err)
			}
		}

		// create the services in the namespace
		for svcNum := 0; svcNum < g.ServerServicesPerNamespace; svcNum++ {
			deployBackend := svcNum % g.ServerDeploymentsPerNamespace
			service := fortio.FortioService{
				Name:                fmt.Sprintf("ns%d-service-%d", nsNum, svcNum),
				Namespace:           "ns" + fmt.Sprint(nsNum),
				TargetPort:          8080,
				ServiceBackendLabel: "fortio-service-" + fmt.Sprint(deployBackend),
			}

			err := yaml.CreateYamlFile(fmt.Sprintf("%s/2-%d-service.yaml", targetDirectory, svcNum), &service)
			if err != nil {
				return fmt.Errorf("failed to create service yaml file: %w", err)
			}
		}

		// create all client deployments
		for clientNum := 0; clientNum < g.ClientDeploymentsPerNamespace; clientNum++ {
			svcNum := clientNum % g.ServerServicesPerNamespace
			svcName := fmt.Sprintf("ns%d-service-%d", nsNum, svcNum)
			appLabel := "fortio-client-" + fmt.Sprint(clientNum)
			requestPort := 8080

			client := fortio.FortioClientDeployment{
				Name:         fmt.Sprintf("ns%d-client-%d", nsNum, clientNum),
				Namespace:    namespace.Name,
				Replicas:     g.ClientReplicasPerDeployment,
				RequestURL:   fmt.Sprintf("ns%d-service-%d", nsNum, svcNum),
				RequestPort:  requestPort,
				AppLabel:     appLabel,
				QPS:          fmt.Sprintf("%d", g.ClientQPS),
				NodeSelector: "scenario: highcount",
			}
			err := yaml.CreateYamlFile(fmt.Sprintf("%s/3-%d-client.yaml", targetDirectory, clientNum), &client)
			if err != nil {
				return fmt.Errorf("failed to create client yaml file: %w", err)
			}

			fqdnPolicy := fortio.ClientFQDNPolicy{
				Name:                   fmt.Sprintf("ns%d-client-%d-fqdn-policy", nsNum, clientNum),
				Namespace:              namespace.Name,
				AppLabel:               appLabel,
				ToPort:                 requestPort,
				TargetServiceName:      svcName,
				TargetServiceNamespace: namespace.Name,
			}
			err = yaml.CreateYamlFile(fmt.Sprintf("%s/4-%d-fqdn-allow.yaml", targetDirectory, clientNum), &fqdnPolicy)
			if err != nil {
				return fmt.Errorf("failed to create client fqdn allow yaml file: %w", err)
			}
		}

		toolbox := fortio.ToolboxDeployment{
			Name:         fmt.Sprintf("ns%d-toolbox", nsNum),
			Namespace:    "ns" + fmt.Sprint(nsNum),
			NodeSelector: "scenario: highcount",
			Replicas:     1,
		}
		err = yaml.CreateYamlFile(fmt.Sprintf("%s/4-toolbox.yaml", targetDirectory), &toolbox)
		if err != nil {
			return fmt.Errorf("failed to create client yaml file: %w", err)
		}

	}

	return nil
}
