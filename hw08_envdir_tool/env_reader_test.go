package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"syscall"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("simple test", func(t *testing.T) {
		envs, err := ReadDir("testdata/env")
		testdataEnvs := Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", false},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}
		require.NoError(t, err)
		require.Equal(t, testdataEnvs, envs)
	})

	t.Run("not a directory", func(t *testing.T) {
		_, err := ReadDir("/dev/null")
		require.ErrorIs(t, err, syscall.ENOTDIR)
	})

	t.Run("no permissions", func(t *testing.T) {
		_, err := ReadDir("/root")
		require.ErrorIs(t, err, os.ErrPermission)
	})
}
