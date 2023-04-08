package query

import (
	"context"
	"fmt"
	"k8s_ui/utils"
	"strings"

	v2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func HPAs(ns string, wide bool) ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	header := []string{"NAME", "REFERENCE", "TARGETS", "MINPODS", "MAXPODS", "REPLICAS", "AGE"}
	if wide {
		header = append(header, "CONDITIONS")
	}

	hpas, err := client.AutoscalingV2().HorizontalPodAutoscalers(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data := [][]string{header}

	for _, hpa := range hpas.Items {
		row := make([]string, len(header))
		row[0] = hpa.Name
		row[1] = hpa_reference(&hpa)
		row[2] = hpa_targets(&hpa)
		row[3] = hpa_min_pods(&hpa)
		row[4] = hpa_max_pods(&hpa)
		row[5] = hpa_replicas(&hpa)
		row[6] = utils.HumanElapsedTime(hpa.CreationTimestamp.Time)
		if wide {
			row[7] = hpa_conditions(&hpa)
		}
		data = append(data, row)
	}

	return data, nil
}

func hpa_reference(hpa *v2.HorizontalPodAutoscaler) string {
	return fmt.Sprintf("%s/%s", hpa.Spec.ScaleTargetRef.Kind, hpa.Spec.ScaleTargetRef.Name)
}

func hpa_targets(hpa *v2.HorizontalPodAutoscaler) string {
	var targets []string
	for _, metric := range hpa.Spec.Metrics {
		switch metric.Type {
		// example: hits-per-second on an Ingress object
		case v2.ObjectMetricSourceType:
			if metric.Object.Target.AverageUtilization != nil {
				targets = append(targets, fmt.Sprintf("Obj(%s/%d)", metric.Object.Metric.Name, *metric.Object.Target.AverageUtilization))
			}

		// example: transactions-processed-per-second
		case v2.PodsMetricSourceType:
			if metric.Pods.Target.AverageUtilization != nil {
				targets = append(targets, fmt.Sprintf("Pod(%s:%d)", metric.Pods.Metric.Name, *metric.Pods.Target.AverageUtilization))
			}

		// example: CPU or memory usage
		case v2.ResourceMetricSourceType:
			if metric.Resource.Target.AverageUtilization != nil {
				targets = append(targets, fmt.Sprintf("Res(%s:%d%%)", metric.Resource.Name, *metric.Resource.Target.AverageUtilization))
			}

		// example: CPU or memory usage
		case v2.ContainerResourceMetricSourceType:
			// TODO: test this somehow? :)
			targets = append(targets, fmt.Sprintf("CRes: %s %s", metric.ContainerResource.Name, metric.ContainerResource.Target.String()))

		case v2.ExternalMetricSourceType:
			//TODO: test this somehow? :)
			targets = append(targets, fmt.Sprintf("Ext %s %s", metric.External.Metric.Name, metric.External.Target.String()))
		}
	}

	return strings.Join(targets, ",")
}

func hpa_min_pods(hpa *v2.HorizontalPodAutoscaler) string {
	if hpa.Spec.MinReplicas == nil {
		return "1" // default is 1
	}

	return fmt.Sprintf("%d", *hpa.Spec.MinReplicas)
}

func hpa_max_pods(hpa *v2.HorizontalPodAutoscaler) string {
	return fmt.Sprintf("%d", hpa.Spec.MaxReplicas)
}

func hpa_replicas(hpa *v2.HorizontalPodAutoscaler) string {
	return fmt.Sprintf("%d", hpa.Status.CurrentReplicas)
}

func hpa_conditions(hpa *v2.HorizontalPodAutoscaler) string {
	var conditions []string
	for _, condition := range hpa.Status.Conditions {
		conditions = append(conditions, fmt.Sprintf("%s:%s", condition.Type, condition.Status))
	}
	return strings.Join(conditions, ",")
}
