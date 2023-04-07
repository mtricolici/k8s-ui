package k8s

import (
	"k8s_ui/k8s/query"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient struct {
	client *kubernetes.Clientset
}

func NewKubernetesClient() *KubernetesClient {
	return &KubernetesClient{
		client: nil,
	}
}

func (c *KubernetesClient) initClient() error {
	if c.client == nil {
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return err
		}

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			return err
		}

		c.client = client
	}

	return nil
}

func (c *KubernetesClient) GetNamespaces() ([][]string, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return query.Namespaces(c.client)
}
