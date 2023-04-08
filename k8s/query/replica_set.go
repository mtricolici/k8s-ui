package query

import (
	"context"
	"fmt"
	"k8s_ui/utils"
	"strings"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ReplicaSets(ns string, wide bool) ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	header := []string{"NAME", "DESIRED", "CURRENT", "READY", "AGE"}
	if wide {
		header = append(header, "CONTAINERS", "IMAGES", "SELECTOR")
	}

	replicaSets, err := client.AppsV1().ReplicaSets(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data := [][]string{header}

	for _, rs := range replicaSets.Items {
		row := make([]string, len(header))
		row[0] = rs.Name
		row[1] = rs_desired(&rs)
		row[2] = fmt.Sprintf("%d", rs.Status.Replicas)
		row[3] = fmt.Sprintf("%d", rs.Status.ReadyReplicas)
		row[4] = utils.HumanElapsedTime(rs.CreationTimestamp.Time)
		if wide {
			row[5] = rs_containers(&rs)
			row[6] = rs_images(&rs)
			row[7] = rs_selector(&rs)
		}
		data = append(data, row)
	}

	return data, nil
}

func rs_desired(rs *v1.ReplicaSet) string {
	if rs.Spec.Replicas == nil {
		return "1" // default is 1
	}

	return fmt.Sprintf("%d", *rs.Spec.Replicas)
}

func rs_containers(rs *v1.ReplicaSet) string {
	var data []string
	for _, container := range rs.Spec.Template.Spec.Containers {
		data = append(data, container.Name)
	}
	return strings.Join(data, ",")
}

func rs_images(rs *v1.ReplicaSet) string {
	var data []string
	for _, container := range rs.Spec.Template.Spec.Containers {
		data = append(data, container.Image)
	}
	return strings.Join(data, ",")
}

func rs_selector(rs *v1.ReplicaSet) string {
	if rs.Spec.Selector == nil {
		return "<none>"
	}

	return utils.SelectorToString(rs.Spec.Selector.MatchLabels)
}
