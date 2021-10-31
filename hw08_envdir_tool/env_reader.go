package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

const (
	newLine   = 0x0a
	emptyByte = 0x00
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadFileLine(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileReader := bufio.NewReader(file)
	isPrefix := true
	value := make([]byte, 0)
	for isPrefix {
		line, finishedReading, err := fileReader.ReadLine()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		value = append(value, line...)
		isPrefix = finishedReading
	}

	return bytes.ReplaceAll(value, []byte{emptyByte}, []byte{newLine}), nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	directory, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer directory.Close()
	dirFiles, err := directory.ReadDir(0)
	if err != nil {
		return nil, err
	}
	resultEnv := make(Environment)
	for _, dirFile := range dirFiles {
		if dirFile.IsDir() {
			continue
		}
		filePath := path.Join(dir, dirFile.Name())
		fileLine, err := ReadFileLine(filePath)
		if err != nil {
			continue
		}
		firstLine := strings.TrimRight(string(fileLine), ` \t`)
		envValue := EnvValue{
			Value:      firstLine,
			NeedRemove: false,
		}
		if firstLine == "" {
			envValue.NeedRemove = true
		}
		resultEnv[dirFile.Name()] = envValue
	}

	return resultEnv, nil
}
