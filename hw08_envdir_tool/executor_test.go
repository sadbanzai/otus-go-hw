package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("exit code", func(t *testing.T) {
		exitCode := RunCmd([]string{"sh", "-c", "exit 2"}, Environment{})
		require.Equal(t, 2, exitCode)
	})

	t.Run("exit code with env", func(t *testing.T) {
		exitCode := RunCmd([]string{"sh", "-c", "exit $EXITCODE"}, Environment{"EXITCODE": EnvValue{Value: "2", NeedRemove: false}})
		require.Equal(t, 2, exitCode)
	})

	t.Run("unknown error exit code", func(t *testing.T) {
		exitCode := RunCmd([]string{"abcd"}, Environment{})
		require.Equal(t, UnknownErrorExitCode, exitCode)
	})
}
