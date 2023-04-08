package query

import (
	"context"
	"fmt"
	"k8s_ui/utils"
	"strings"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DaemonSets(ns string, wide bool) ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	header := []string{"NAME", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE", "AGE"}
	if wide {
		header = append(header, "NODE-SELECTOR", "CONTAINERS", "IMAGES")
	}

	daemonSets, err := client.AppsV1().DaemonSets(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data := [][]string{header}

	for _, ds := range daemonSets.Items {
		row := make([]string, len(header))
		row[0] = ds.Name
		row[1] = ds_desired(&ds)
		row[2] = ds_current(&ds)
		row[3] = ds_ready(&ds)
		row[4] = ds_up_to_date(&ds)
		row[5] = ds_available(&ds)
		row[6] = utils.HumanElapsedTime(ds.CreationTimestamp.Time)
		if wide {
			row[7] = ds_node_selector(&ds)
			row[8] = ds_containers(&ds)
			row[9] = ds_images(&ds)
		}
		data = append(data, row)
	}

	return data, nil
}

func ds_desired(ds *v1.DaemonSet) string {
	return fmt.Sprintf("%d", ds.Status.DesiredNumberScheduled)
}

func ds_current(ds *v1.DaemonSet) string {
	return fmt.Sprintf("%d", ds.Status.CurrentNumberScheduled)
}

func ds_ready(ds *v1.DaemonSet) string {
	return fmt.Sprintf("%d", ds.Status.NumberReady)
}

func ds_up_to_date(ds *v1.DaemonSet) string {
	return fmt.Sprintf("%d", ds.Status.UpdatedNumberScheduled)
}

func ds_available(ds *v1.DaemonSet) string {
	return fmt.Sprintf("%d", ds.Status.NumberAvailable)
}

func ds_node_selector(ds *v1.DaemonSet) string {
	if ds.Spec.Template.Spec.NodeSelector == nil {
		return ""
	}
	return fmt.Sprintf("%v", ds.Spec.Template.Spec.NodeSelector)
}

func ds_containers(ds *v1.DaemonSet) string {
	var containers []string
	for _, container := range ds.Spec.Template.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return strings.Join(containers, ",")
}

func ds_images(ds *v1.DaemonSet) string {
	var images []string
	for _, container := range ds.Spec.Template.Spec.Containers {
		images = append(images, container.Image)
	}
	return strings.Join(images, ",")
}
