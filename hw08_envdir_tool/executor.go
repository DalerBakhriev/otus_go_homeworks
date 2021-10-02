package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func getModifiedEnvVars(env Environment) []string {
	envVars := os.Environ()
	currentEnvVars := make(map[string]string)
	for _, envVar := range envVars {
		nameAndValue := strings.SplitN(envVar, "=", 2)
		name, value := nameAndValue[0], nameAndValue[1]
		currentEnvVars[name] = value
	}
	for name, newValue := range env {
		if newValue.NeedRemove {
			delete(currentEnvVars, name)
			continue
		}
		currentEnvVars[name] = newValue.Value
	}
	modifiedEnvVars := make([]string, 0, len(currentEnvVars))
	for name, value := range currentEnvVars {
		modifiedEnvVars = append(modifiedEnvVars, strings.Join([]string{name, value}, "="))
	}
	return modifiedEnvVars
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	appName, args := cmd[0], cmd[1:]
	command := exec.Command(appName, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	modifiedEnvVars := getModifiedEnvVars(env)
	command.Env = modifiedEnvVars

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			returnCode = exitError.ExitCode()
		}
	}

	return returnCode
}
