package shell

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func Command(cmd string) (string, error) {
	var args = append([]string{"cmd", "/C"}, strings.Fields(strings.TrimSpace(cmd))...)

	command := exec.Command(args[0], args[1:]...)
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("executing command error: %w", err)
	}
	return string(output), nil
}
