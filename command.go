package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type Command struct {
	stdInFunction func()
	stdInDuration int
	stdin         io.Reader
	stdout        io.Writer
	stderr        io.Writer
}

func (c *Command) Run(cmd string) ([]byte, error) {
	return exec.Command("bash", cmd).Output()
}

func (c *Command) RunWithoutBash(cmd string) ([]byte, error) {
	return exec.Command(cmd).Output()
}

// RunCommand runs a command with custom stdin, stdout, stderr. If any is nil, uses os.Stdin, os.Stdout, os.Stderr respectively.
func (t *Command) RunCommand(path string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	if path != "" {
		cmd.Dir = path
	}
	// stdin
	if t.stdin != nil {
		cmd.Stdin = t.stdin
	} else {
		cmd.Stdin = os.Stdin
	}
	// stdout
	if t.stdout != nil {
		cmd.Stdout = t.stdout
	} else {
		cmd.Stdout = os.Stdout
	}
	// stderr
	if t.stderr != nil {
		cmd.Stderr = t.stderr
	} else {
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

// RunCommandStd eski davranış için: sadece path, name, arg alır ve çıktıyı ekrana basar
// stdin setter
func (t *Command) SetStdin(r io.Reader) {
	t.stdin = r
}

// stdout setter
func (t *Command) SetStdout(w io.Writer) {
	t.stdout = w
}

// stderr setter
func (t *Command) SetStderr(w io.Writer) {
	t.stderr = w
}
func (t *Command) RunWithPipe(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	go func() {
		if t.stdInFunction != nil {
			time.Sleep(time.Duration(t.stdInDuration) * time.Second)
			t.stdInFunction()
		}
	}()

	err := cmd.Run()
	fmt.Println(err)
}

func (t *Command) AddStdIn(duration int, f func()) {
	t.stdInFunction = f
	t.stdInDuration = duration
}
