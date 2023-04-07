package query

import (
	"context"
	"fmt"
	"k8s_ui/utils"
	"strings"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Ingresses(ns string, wide bool) ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	header := []string{"NAME", "CLASS", "HOSTS", "ADDRESSES", "PORTS", "AGE"}
	if wide {
		header = append(header, "TLS", "RULES")
	}

	ingresses, err := client.NetworkingV1().Ingresses(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data := [][]string{header}

	for _, ing := range ingresses.Items {
		row := make([]string, len(header))
		row[0] = ing.Name
		row[1] = ingress_class(&ing)
		row[2] = ingress_hosts(&ing)
		row[3] = ingress_addresses(&ing)
		row[4] = ingress_ports(&ing)
		row[5] = utils.HumanElapsedTime(ing.CreationTimestamp.Time)
		if wide {
			row[6] = ingress_tls(&ing)
			row[7] = ingress_rules(&ing)
		}
		data = append(data, row)
	}

	return data, nil
}

func ingress_class(i *v1.Ingress) string {
	if i.Spec.IngressClassName == nil {
		return "<none>"
	}

	return *i.Spec.IngressClassName
}

func ingress_hosts(i *v1.Ingress) string {
	if len(i.Spec.Rules) == 0 {
		return "*"
	}

	var hosts []string
	for _, rule := range i.Spec.Rules {
		if len(rule.Host) > 0 {
			hosts = append(hosts, rule.Host)
		}

	}
	if len(hosts) > 0 {
		return strings.Join(hosts, ",")
	}

	return "*"
}

func ingress_addresses(i *v1.Ingress) string {
	var addresses []string
	for _, lb := range i.Status.LoadBalancer.Ingress {
		if lb.IP != "" {
			addresses = append(addresses, lb.IP)
		} else if lb.Hostname != "" {
			addresses = append(addresses, lb.Hostname)
		}
	}
	return strings.Join(addresses, ",")
}

func ingress_ports(i *v1.Ingress) string {
	var ports []string
	for _, rule := range i.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service != nil {
				ports = append(ports, fmt.Sprintf("%d", path.Backend.Service.Port.Number))
			}
		}
	}
	return strings.Join(ports, ",")
}

func ingress_tls(i *v1.Ingress) string {
	var tls []string
	for _, tlsSpec := range i.Spec.TLS {
		tls = append(tls, tlsSpec.SecretName)
	}
	return strings.Join(tls, ",")
}

func ingress_rules(i *v1.Ingress) string {
	var rules []string
	for _, rule := range i.Spec.Rules {
		if rule.IngressRuleValue.HTTP != nil {
			for _, path := range (*rule.IngressRuleValue.HTTP).Paths {
				rules = append(rules, rule.Host+path.Path)
			}
		}
	}
	return strings.Join(rules, ",")
}
