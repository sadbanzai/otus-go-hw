package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("exit code", func(t *testing.T) {
		exitCode := RunCmd([]string{"sh", "-c", "exit 2"}, Environment{})
		require.Equal(t, 2, exitCode)
	})

	t.Run("exit code with env", func(t *testing.T) {
		env := Environment{
			"EXITCODE": EnvValue{Value: "2", NeedRemove: false},
		}
		exitCode := RunCmd([]string{"sh", "-c", "exit $EXITCODE"}, env)
		require.Equal(t, 2, exitCode)
	})

	t.Run("unknown error exit code", func(t *testing.T) {
		exitCode := RunCmd([]string{"abcd"}, Environment{})
		require.Equal(t, UnknownErrorExitCode, exitCode)
	})
}
