package query

import (
	"context"
	"fmt"
	"k8s_ui/utils"
	"strings"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StatefulSets(ns string, wide bool) ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	header := []string{"NAME", "READY", "CURRENT", "UPDATED", "AVAILABLE", "AGE"}
	if wide {
		header = append(header, "CONTAINERS", "IMAGES", "VOLUMES")
	}

	statefulSets, err := client.AppsV1().StatefulSets(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data := [][]string{header}

	for _, ss := range statefulSets.Items {
		row := make([]string, len(header))
		row[0] = ss.Name
		row[1] = statefulset_ready(&ss)
		row[2] = fmt.Sprintf("%d", ss.Status.CurrentReplicas)
		row[3] = fmt.Sprintf("%d", ss.Status.UpdatedReplicas)
		row[4] = fmt.Sprintf("%d", ss.Status.AvailableReplicas)
		row[5] = utils.HumanElapsedTime(ss.CreationTimestamp.Time)
		if wide {
			row[6] = statefulset_containers(&ss)
			row[7] = statefulset_images(&ss)
			row[8] = statefulset_volumes(&ss)
		}
		data = append(data, row)
	}

	return data, nil
}

func statefulset_ready(ss *v1.StatefulSet) string {
	return fmt.Sprintf("%d/%d", ss.Status.ReadyReplicas, ss.Status.Replicas)
}

func statefulset_containers(ss *v1.StatefulSet) string {
	var names []string
	for _, container := range ss.Spec.Template.Spec.Containers {
		names = append(names, container.Name)
	}
	return strings.Join(names, ",")
}

func statefulset_images(ss *v1.StatefulSet) string {
	var images []string
	for _, container := range ss.Spec.Template.Spec.Containers {
		images = append(images, container.Image)
	}
	return strings.Join(images, ",")
}

func statefulset_volumes(ss *v1.StatefulSet) string {
	var volumes []string
	for _, volume := range ss.Spec.Template.Spec.Volumes {
		volumes = append(volumes, volume.Name)
	}

	return strings.Join(volumes, ",")
}
