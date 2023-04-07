package query

import (
	"context"
	"k8s_ui/utils"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PersistentVolumeClaims(ns string, wide bool) ([][]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	header := []string{"NAME", "STATUS", "VOLUME", "CAPACITY", "ACCESS-MODES", "STORAGE-CLASS", "AGE"}
	if wide {
		header = append(header, "VOLUME-MODE")
	}

	pvcs, err := client.CoreV1().PersistentVolumeClaims(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	data := [][]string{header}

	for _, pvc := range pvcs.Items {
		row := make([]string, len(header))
		row[0] = pvc.Name
		row[1] = string(pvc.Status.Phase)
		row[2] = pvc.Spec.VolumeName
		row[3] = pvc_capacity(&pvc)
		row[4] = pvc_access_mods(&pvc)
		row[5] = pvc_storage_class(&pvc)
		row[6] = utils.HumanElapsedTime(pvc.CreationTimestamp.Time)
		if wide {
			row[7] = pvc_volume_mode(&pvc)
		}
		data = append(data, row)
	}

	return data, nil
}

func pvc_capacity(pvc *v1.PersistentVolumeClaim) string {
	capacity, ok := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	if !ok {
		return "0"
	}

	return capacity.String()
}

func pvc_access_mods(pvc *v1.PersistentVolumeClaim) string {
	var modes []string
	for _, mode := range pvc.Spec.AccessModes {
		modes = append(modes, string(mode))
	}
	return strings.Join(modes, ",")
}

func pvc_storage_class(pvc *v1.PersistentVolumeClaim) string {
	if pvc.Spec.StorageClassName == nil {
		return ""
	}
	return *pvc.Spec.StorageClassName
}

func pvc_volume_mode(pvc *v1.PersistentVolumeClaim) string {
	if pvc.Spec.VolumeMode == nil {
		return ""
	}

	return string(*pvc.Spec.VolumeMode)
}
