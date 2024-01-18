package shell

import (
	"bytes"
	"os/exec"
)

type Shell interface {
	// Shell return at End all
	Shell(command ...string) (stdout, stderr string, err error)

	// StreamShell send per line
	StreamShell(stdout, stderr chan string, command ...string) (err error)

	// StreamCallbackShell call per line
	StreamCallbackShell(stdout, stderr func(line string), command ...string) (err error)
}

func Exec(command ...string) (stdout, stderr string, err error) {
	return DefaultShell(command...)
}

func run(cmd *exec.Cmd) (stdout, stderr string, err error) {
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()
	stdout = outb.String()
	stderr = errb.String()
	return
}
