package query

import (
	"context"
	"fmt"
	"k8s_ui/utils"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Pods(ns string, wide bool) ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	header := []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE"}
	if wide {
		header = append(header, "IP", "NODE")
	}

	pods, err := client.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data := [][]string{header}

	for _, pod := range pods.Items {
		row := make([]string, len(header))
		row[0] = pod.Name
		row[1] = pod_ready(&pod)
		row[2] = pod_status(&pod)
		row[3] = pod_restarts(&pod)
		row[4] = utils.HumanElapsedTime(pod.CreationTimestamp.Time)
		if wide {
			row[5] = pod_ip(&pod)
			row[6] = pod_node(&pod)
		}
		data = append(data, row)
	}

	return data, nil
}

func pod_ready(pod *v1.Pod) string {
	total := len(pod.Spec.Containers)
	ready := 0

	for _, status := range pod.Status.ContainerStatuses {
		if status.Ready {
			ready++
		}
	}

	return fmt.Sprintf("%d/%d", ready, total)
}

func pod_status(pod *v1.Pod) string {
	return string(pod.Status.Phase)
}

func pod_restarts(pod *v1.Pod) string {
	var restarts int32

	if len(pod.Status.ContainerStatuses) > 0 {
		restarts = pod.Status.ContainerStatuses[0].RestartCount
	}

	return fmt.Sprintf("%d", restarts)
}

func pod_ip(pod *v1.Pod) string {
	if len(pod.Status.PodIP) > 0 {
		return pod.Status.PodIP
	}
	return "<none>"
}

func pod_node(pod *v1.Pod) string {
	if len(pod.Spec.NodeName) > 0 {
		return pod.Spec.NodeName
	}
	return "<none>"
}
