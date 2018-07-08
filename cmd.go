package kiosk

import (
	"io"
	"os"
	"os/exec"
)

var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr

//spawn returns the output from the program supplied
func spawnCmd(prog string, args ...string) (string, error) {
	out, err := exec.Command(prog, args...).Output()
	if err != nil {
		return "", err
	}
	return (string)(out), nil
}

func runCmd(prog string, args ...string) error {
	cmd := exec.Command(prog, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
