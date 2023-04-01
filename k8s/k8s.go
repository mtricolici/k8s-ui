package k8s

import (
	"k8s_ui/utils"
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

func (client *K8SClient) GetPods(ns string, wide bool) ([][]string, error) {
	args := client.get_pods
	args = append(args, "-n", ns)
	if wide {
		args = append(args, "-o", "wide")
	}
	return client.exec(args)
}

func (client *K8SClient) GetResources(ns, resource string, wide bool) ([][]string, error) {
	args := []string{
		"get", resource, "-n", ns, "--sort-by", ".metadata.name",
	}

	if wide {
		args = append(args, "-o", "wide")
	}
	return client.exec(args)
}

func (client *K8SClient) GetPodContainerNames(ns, pod string) ([]string, error) {
	args := []string{
		"get",
		"pod",
		pod,
		"-n",
		ns,
		"-o",
		`jsonpath='{.spec.containers[*].name}'`,
	}

	result, err := client.exec(args)
	if err != nil {
		return nil, err
	}
	// I don't know why ' appears here :D wtf. Let's just filter them
	containers := []string{}
	for i := range result[0] {
		containers = append(containers, utils.ReplaceSpecialChars(result[0][i]))
	}

	return containers, nil
}
