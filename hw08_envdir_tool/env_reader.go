package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrTooLongLine = errors.New("too long line")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	envs := make(Environment, len(files))
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		hasEqualSymbol := strings.Contains(info.Name(), "=")
		if hasEqualSymbol {
			continue
		}
		if info.Size() == 0 {
			envs[info.Name()] = EnvValue{"", true}
			continue
		}
		filePtr, err := os.Open(filepath.Join(dir, info.Name()))
		if err != nil {
			return nil, err
		}
		br := bufio.NewReader(filePtr)
		line, _, err := br.ReadLine()
		if err != nil {
			return nil, err
		}
		line = bytes.Replace(line, []byte("\x00"), []byte("\n"), -1)
		line = bytes.TrimRight(line, " \t")
		envs[info.Name()] = EnvValue{string(line), false}
	}
	return envs, nil
}
