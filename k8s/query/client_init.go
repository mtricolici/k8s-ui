package query

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	_client *kubernetes.Clientset
)

func getClient() (*kubernetes.Clientset, error) {
	if _client == nil {
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, err
		}

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		_client = client
	}

	return _client, nil
}
