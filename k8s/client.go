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
	switch resource {
	case "pod":
		return query.Pods(ns, wide)
	case "service":
		return query.Services(ns, wide)
	case "deployment":
		return query.Deployments(ns, wide)
	case "ingress":
	case "pvc":
	case "daemonset":
	case "replicaset":
	case "statefulset":
	case "all":
	}

	return nil, fmt.Errorf("get resource '%s' not implemented yet", resource)
}
