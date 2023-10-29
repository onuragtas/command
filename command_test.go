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
		StdOutWriter: writeToStdout,
		StdErrWriter: writeToStderr,
		StdInData:    "Bu komuta giriş verisi gönderilecek.",
	}

	cmd.RunCommand("./", "ls", "-la")
	//cmd.RunCommand("./", "watch", "-n", "2", "ls -la")
}
