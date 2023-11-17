package query

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	_client *kubernetes.Clientset
)

func getClient() (*kubernetes.Clientset, error) {
	if _client == nil {

		kubeConfigPath := os.Getenv("KUBECONFIG")
		if kubeConfigPath == "" || !fileExists(kubeConfigPath) {
			// Use default ~/.kube/config
			kubeConfigPath = clientcmd.RecommendedHomeFile
		}

		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
