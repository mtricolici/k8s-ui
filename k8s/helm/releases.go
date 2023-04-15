package helm

import (
	"errors"
	"fmt"
	"k8s_ui/utils"

	"k8s.io/apimachinery/pkg/util/yaml"
)

func Releases(ns string) ([][]string, error) {
	output, err := utils.Execute("helm ls -a -n %s -o yaml", ns)
	if err != nil {
		return nil, err
	}

	return parse_releases(output)
}

func parse_releases(output []byte) ([][]string, error) {
	var raw_data interface{}
	err := yaml.Unmarshal(output, &raw_data)
	if err != nil {
		return nil, err
	}

	// expecting to be an array
	if data, ok := raw_data.([]any); ok {
		result := [][]string{
			{"NAME", "REVISION", "UPDATED", "STATUS", "CHART", "APP VERSION"},
		}

		for _, raw_item := range data {
			// expecting to have a map of key:value
			if item, ok := raw_item.(map[string]any); ok {
				raw := []string{
					fmt.Sprintf("%s", item["name"]),
					fmt.Sprintf("%s", item["revision"]),
					fmt.Sprintf("%s", item["updated"]),
					fmt.Sprintf("%s", item["status"]),
					fmt.Sprintf("%s", item["chart"]),
					fmt.Sprintf("%s", item["app_version"]),
				}
				result = append(result, raw)
			} else {
				goto releases_error
			}
		}

		return result, nil
	}

releases_error:
	return nil, errors.New("bad yaml output")
}
