package query

import (
	"context"
	"fmt"
	"k8s_ui/utils"
	"strings"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Deployments(ns string, wide bool) ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	header := []string{"NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE"}
	if wide {
		header = append(header, "CONTAINERS", "IMAGES", "SELECTOR")
	}

	deployments, err := client.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data := [][]string{header}

	for _, deployment := range deployments.Items {
		row := make([]string, len(header))
		row[0] = deployment.Name
		row[1] = deployment_ready(deployment)
		row[2] = deployment_up_to_date(deployment)
		row[3] = deployment_available(deployment)
		row[4] = utils.HumanElapsedTime(deployment.CreationTimestamp.Time)
		if wide {
			row[5] = deployment_containers(deployment)
			row[6] = deployment_images(deployment)
			row[7] = deployment_selectors(deployment)
		}
		data = append(data, row)
	}

	return data, nil
}

func deployment_ready(d v1.Deployment) string {
	return fmt.Sprintf("%d/%d", d.Status.ReadyReplicas, d.Status.Replicas)
}

func deployment_up_to_date(d v1.Deployment) string {
	return fmt.Sprintf("%d", d.Status.UpdatedReplicas)
}

func deployment_available(d v1.Deployment) string {
	return fmt.Sprintf("%d", d.Status.AvailableReplicas)
}

func deployment_containers(d v1.Deployment) string {
	names := []string{}

	for _, container := range d.Spec.Template.Spec.Containers {
		names = append(names, container.Name)
	}

	return strings.Join(names, ",")
}

func deployment_images(d v1.Deployment) string {
	images := []string{}

	for _, container := range d.Spec.Template.Spec.Containers {
		images = append(images, container.Image)
	}

	return strings.Join(images, ",")
}

func deployment_selectors(d v1.Deployment) string {
	if d.Spec.Selector == nil {
		return ""
	}

	selectors := []string{}
	for k, v := range (*d.Spec.Selector).MatchLabels {
		selectors = append(selectors, k+"="+v)
	}

	return strings.Join(selectors, ",")
}
