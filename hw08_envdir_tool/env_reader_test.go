package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	currDir, err := os.Getwd()
	assert.NoError(t, err)
	dirPath := path.Join(currDir, "testdata", "env")

	env, err := ReadDir(dirPath)
	assert.NoError(t, err)

	correctMapRes := map[string]EnvValue{
		"BAR":   {Value: "bar", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: true},
		"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": {Value: "\"hello\"", NeedRemove: false},
		"UNSET": {Value: "", NeedRemove: true},
	}
	expectedEnvRes := Environment(correctMapRes)
	assert.EqualValues(t, env, expectedEnvRes)
}
