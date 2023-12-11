package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", "out.txt", 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("no permissions from-file", func(t *testing.T) {
		err := Copy("/etc/sudoers", "out.txt", 0, 0)
		require.ErrorIs(t, err, os.ErrPermission)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 10000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("limit exceeds file size", func(t *testing.T) {
		toFile, _ := os.CreateTemp("/tmp", "out")
		err := Copy("testdata/input.txt", toFile.Name(), 0, 10000)
		require.NoError(t, err)
	})
}
