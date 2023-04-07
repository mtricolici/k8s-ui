package k8s

import (
	"fmt"
	"k8s_ui/k8s/query"
)

func GetNamespaces() ([][]string, error) {
	return query.Namespaces()
}

func GetPodContainerNames(ns, pod string) ([]string, error) {
	return query.Pod_containers(ns, pod)
}

func GetResources(ns, resource string, wide bool) ([][]string, error) {
	return nil, fmt.Errorf("get resource '%s' not implemented yet", resource)
}
