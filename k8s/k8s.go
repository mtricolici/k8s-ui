package k8s

import (
	"os/exec"
	"regexp"
	"strings"
)

const (
	get_namespaces = "get ns --sort-by .metadata.name"
	get_pods       = "get po --sort-by .metadata.name"
)

type K8SClient struct {
	get_namespaces []string
	get_pods       []string
}

func NewK8SClient() *K8SClient {
	client := K8SClient{
		get_namespaces: strings.Split(get_namespaces, " "),
		get_pods:       strings.Split(get_pods, " "),
	}
	return &client
}

func (client *K8SClient) exec(args []string) ([][]string, error) {
	out, err := exec.Command("kubectl", args[:]...).Output()
	if err != nil {
		return nil, err
	}

	result := [][]string{}

	for _, line := range strings.Split(strings.TrimSuffix(string(out), "\n"), "\n") {

		items := regexp.MustCompile(`\s+`).Split(line, -1)
		result = append(result, items)
	}

	return result[:], nil
}

func (client *K8SClient) GetNamespaces() ([][]string, error) {
	return client.exec(client.get_namespaces)
}

func (client *K8SClient) GetPods(ns string) ([][]string, error) {
	args := client.get_pods
	args = append(args, "-n", ns)
	return client.exec(args)
}
