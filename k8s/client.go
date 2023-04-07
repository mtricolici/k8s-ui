package k8s

import (
	"k8s_ui/k8s/query"
)

func GetNamespaces() ([][]string, error) {
	return query.Namespaces()
}

func GetPodContainerNames(ns, pod string) ([]string, error) {
	return query.Pod_containers(ns, pod)
}
