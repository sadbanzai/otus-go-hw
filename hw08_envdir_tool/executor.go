package main

import (
	"os"
	"os/exec"
)

var (
	// UnknownErrorExitCode returns when cannot run child
	UnknownErrorExitCode = 111
)

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

	runCmd := exec.Command(cmd[0], cmd[1:]...)
	runCmd.Env = os.Environ()
	runCmd.Stdin = os.Stdin
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	err := runCmd.Run()
	switch err := err.(type) {
	case nil:
		return 0
	case *exec.ExitError:
		return err.ExitCode()
	default:
		return UnknownErrorExitCode
	}
}
