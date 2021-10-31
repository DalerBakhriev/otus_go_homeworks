package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	cmd := []string{
		"/bin/bash",
		"./testdata/echo.sh",
		"arg1=1",
		"arg2=2",
	}
	env := Environment(map[string]EnvValue{
		"BAR":   {Value: "bar", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: true},
		"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": {Value: "\"hello\"", NeedRemove: false},
		"UNSET": {Value: "", NeedRemove: true},
	})

	t.Run("correct_redurn_code", func(t *testing.T) {
		returnCode := RunCmd(cmd, env)
		assert.Equal(t, 0, returnCode)
	})

	t.Run("correct_output", func(t *testing.T) {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		RunCmd(cmd, env)
		outC := make(chan string)

		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		w.Close()
		os.Stdout = old
		out := <-outC
		expectedOutput := "HELLO is (\"hello\")\nBAR is (bar)\nFOO is (   foo\nwith new line)\nUNSET is ()\nADDED is ()\nEMPTY is ()\narguments are arg1=1 arg2=2\n"
		assert.Equal(t, expectedOutput, out)
	})
}
