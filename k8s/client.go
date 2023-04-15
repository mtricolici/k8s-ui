package k8s

import (
	"fmt"
	"k8s_ui/k8s/helm"
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
	case "Pod":
		return query.Pods(ns, wide)
	case "Service":
		return query.Services(ns, wide)
	case "Deployment":
		return query.Deployments(ns, wide)
	case "Ingress":
		return query.Ingresses(ns, wide)
	case "pvc":
		return query.PersistentVolumeClaims(ns, wide)
	case "DaemonSet":
		return query.DaemonSets(ns, wide)
	case "ReplicaSet":
		return query.ReplicaSets(ns, wide)
	case "StatefulSet":
		return query.StatefulSets(ns, wide)
	case "Endpoint":
	case "HorizontalPodAutoscaler":
		return query.HPAs(ns, wide)
	case "helm":
		return helm.Releases(ns)
	}

	return nil, fmt.Errorf("get resource '%s' not implemented yet", resource)
}
