package query

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Pod_containers(ns, pod string) ([]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	pod_ob, err := client.CoreV1().Pods(ns).Get(context.Background(), pod, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var containers []string
	for _, container := range pod_ob.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}
