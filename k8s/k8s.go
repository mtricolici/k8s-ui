package k8s

import (
	"fmt"
	"log"
	"os/exec"
)

var (
	get_namespaces []string = []string{"get", "ns", "--no-headers", "--sort-by", ".metadata.name"}
)

func K8s_namespaces() {

	out, err := exec.Command("kubectl", get_namespaces[:]...).Output()
	if err != nil {
		log.Fatal("Error fetching namespaces: ", err)
	}

	fmt.Printf("%s\n", out)
}
