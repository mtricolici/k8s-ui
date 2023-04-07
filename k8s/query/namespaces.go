package query

import (
	"context"
	"k8s_ui/utils"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Namespaces() ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	ns_list, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaces := [][]string{{"NAME", "STATUS", "AGE"}}

	for _, ns := range ns_list.Items {
		row := []string{
			ns.Name,
			namespace_status(ns),
			namespace_age(ns),
		}
		namespaces = append(namespaces, row)
	}

	return namespaces, nil
}

func namespace_status(ns v1.Namespace) string {
	return string(ns.Status.Phase)
}

func namespace_age(ns v1.Namespace) string {
	return utils.HumanElapsedTime(ns.CreationTimestamp.Time)
}
