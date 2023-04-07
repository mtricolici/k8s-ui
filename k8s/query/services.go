package query

import (
	"context"
	"fmt"
	"k8s_ui/utils"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Services(ns string, wide bool) ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	var header []string
	if wide {
		header = []string{"NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE", "SELECTOR"}
	} else {
		header = []string{"NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE"}
	}

	services, err := client.CoreV1().Services(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data := [][]string{header}

	for _, svc := range services.Items {
		row := make([]string, len(header))
		row[0] = svc.Name
		row[1] = string(svc.Spec.Type)
		row[2] = svc.Spec.ClusterIP
		row[3] = svc_external_ips(svc)
		row[4] = svc_ports(svc)
		row[5] = utils.HumanElapsedTime(svc.CreationTimestamp.Time)
		if wide {
			row[6] = svc_selectors(svc)
		}
		data = append(data, row)
	}

	return data, nil
}

func svc_external_ips(svc v1.Service) string {
	if len(svc.Spec.ExternalIPs) > 0 {
		return strings.Join(svc.Spec.ExternalIPs, ",")
	}

	return "<none>"
}

func svc_ports(svc v1.Service) string {
	ports := []string{}

	for _, port := range svc.Spec.Ports {
		portStr := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		if port.NodePort != 0 {
			portStr += fmt.Sprintf(":%d", port.NodePort)
		}
		ports = append(ports, portStr)
	}

	return strings.Join(ports, ",")
}

func svc_selectors(svc v1.Service) string {
	selectors := []string{}
	for k, v := range svc.Spec.Selector {
		selectors = append(selectors, k+"="+v)
	}

	return strings.Join(selectors, ",")
}
