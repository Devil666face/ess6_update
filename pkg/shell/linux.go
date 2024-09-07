//go:build linux

package shell

import (
	"fmt"
	"os/exec"
	"strings"
)

func Command(cmd string) (string, error) {
	var args = strings.Fields(strings.TrimSpace(cmd))

	command := exec.Command(args[0], args[1:]...)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("executing command error: %w", err)
	}
	return string(output), nil
}
