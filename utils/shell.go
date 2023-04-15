package utils

import (
	"fmt"
	"os/exec"
)

func Execute(format string, args ...interface{}) ([]byte, error) {
	command := fmt.Sprintf(format, args...)
	cmd := exec.Command("bash", "-c", command)
	return cmd.Output()
}
