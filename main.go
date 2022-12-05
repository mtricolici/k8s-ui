package main

import (
	"fmt"
	"k8s_ui/k8s"
)

func main() {
	fmt.Printf("namespaces:\n")

	ns := k8s.K8s_namespaces()
	fmt.Println(ns)
}
