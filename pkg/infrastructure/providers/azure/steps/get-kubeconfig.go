package azuresteps

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"k8s.io/client-go/tools/clientcmd"
)

const KubeConfigPerms = 0o600

type GetKubeConfig struct {
	ClusterName       string
	ResourceGroupName string
	SubscriptionID    string
}

func (c *GetKubeConfig) Do(context.Context) error {
	kubeconfigfilepath := clientcmd.RecommendedHomeFile

	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return fmt.Errorf("failed to obtain a credential: %w", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory(c.SubscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	res, err := clientFactory.NewManagedClustersClient().ListClusterUserCredentials(ctx, c.ResourceGroupName, c.ClusterName, nil)
	if err != nil {
		return fmt.Errorf("failed to finish the get managed cluster client request: %w", err)
	}

	err = os.WriteFile(kubeconfigfilepath, []byte(res.Kubeconfigs[0].Value), KubeConfigPerms)
	if err != nil {
		return fmt.Errorf("failed to write kubeconfig to file \"%s\": %w", kubeconfigfilepath, err)
	}

	log.Printf("kubeconfig for cluster \"%s\" in resource group \"%s\" written to \"%s\"\n", c.ClusterName, c.ResourceGroupName, kubeconfigfilepath)
	return nil
}
