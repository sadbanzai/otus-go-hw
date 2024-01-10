package main

import (
	"errors"
	"os"
	"os/exec"
)

var UnknownErrorExitCode = 111

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				return UnknownErrorExitCode
			}
		}
		err := os.Setenv(k, v.Value)
		if err != nil {
			return UnknownErrorExitCode
		}
	}

	runCmd := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	runCmd.Env = os.Environ()
	runCmd.Stdin = os.Stdin
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	err := runCmd.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return 1
	}
	return 0
}
