package k8s

import (
	l "k8s_ui/logger"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var (
	get_namespaces = []string{"get", "ns", "--sort-by", ".metadata.name"}[:]
	get_pods       = []string{"get", "po", "--sort-by", ".metadata.name"}[:]
)

func exec_kubectl(args []string) ([][]string, error) {
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

func exec_get_namespaces() [][]string {
	defer l.LogExecutedTime("exec_get_namespaces")()

	result, err := exec_kubectl(get_namespaces)

	if err != nil {
		log.Panic("Error fetching namespaces: ", err)
	}

	return result
}

func exec_get_pods(namespace string) [][]string {
	defer l.LogExecutedTime("exec_get_pods")()

	result, err := exec_kubectl(append(get_pods, []string{"-n", namespace}...))

	if err != nil {
		log.Panic("Error fetching pods: ", err)
	}

	return result
}
