package command

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	writeToStdout := func(output []byte) {
		fmt.Printf("Custom stdout:\n%s\n", string(output))
	}

	writeToStderr := func(output []byte) {
		fmt.Printf("Custom stderr:\n%s\n", string(output))
	}

	// Command yapısını oluşturun
	cmd := Command{
		StdOutWriter:  writeToStdout,
		StdErrWriter:  writeToStderr,
		Sleep:         1000,
		OutputAndQuit: true,
	}

	cmd.RunCommand("./", "ps", "aux")
	//cmd.RunCommand("./", "ls", "-la")
	//cmd.RunCommand("./", "watch", "-n", "2", "ls -la")
}
