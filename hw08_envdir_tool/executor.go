package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

func modifyEnvVars(env Environment) error {
	// Setting and unsetting env variables
	for name, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(name)
			if err != nil {
				return err
			}
		} else {
			err := os.Setenv(name, value.Value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	err := modifyEnvVars(env)
	if err != nil {
		log.Fatal(err)
	}
	appName, args := cmd[0], cmd[1:]
	command := exec.Command(appName, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			returnCode = exitError.ExitCode()
		}
	}

	return returnCode
}
