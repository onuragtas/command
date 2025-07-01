package command

import (
	"bytes"
	"strings"
	"testing"
)

// Test Run fonksiyonu
func TestRun(t *testing.T) {
	cmd := &Command{}
	out, err := cmd.Run("echo hello")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if !strings.Contains(string(out), "hello") {
		t.Errorf("unexpected output: %s", out)
	}
}

// Test RunWithoutBash fonksiyonu
func TestRunWithoutBash(t *testing.T) {
	cmd := &Command{}
	out, err := cmd.RunWithoutBash("ls")
	if err != nil {
		t.Fatalf("RunWithoutBash failed: %v", err)
	}
	if len(out) == 0 {
		t.Errorf("expected output, got none")
	}
}

// Test RunCommand ve setter fonksiyonları
func TestRunCommandWithSetters(t *testing.T) {
	cmd := &Command{}
	var outBuf, errBuf bytes.Buffer

	cmd.SetStdout(&outBuf)
	cmd.SetStderr(&errBuf)
	cmd.SetStdin(strings.NewReader(""))

	err := cmd.RunCommand("", "echo", "test")
	if err != nil {
		t.Fatalf("RunCommand failed: %v", err)
	}
	if !strings.Contains(outBuf.String(), "test") {
		t.Errorf("unexpected stdout: %s", outBuf.String())
	}
}

// Test RunWithPipe fonksiyonu
func TestRunWithPipe(t *testing.T) {
	cmd := &Command{}
	// Bu fonksiyonun çıktısı doğrudan os.Stdout'a gider, testte gözle kontrol edilebilir.
	cmd.RunWithPipe("echo", "pipe test")
}

// Test AddStdIn fonksiyonu
func TestAddStdIn(t *testing.T) {
	cmd := &Command{}
	called := false
	cmd.AddStdIn(0, func() { 
		called = true
	 })
	cmd.RunWithPipe("echo", "test")
	if !called {
		t.Errorf("AddStdIn fonksiyonu çağrılmadı")
	}
}

func Test(t *testing.T) {

	// Command yapısını oluşturun
	cmd := Command{
	}

	cmd.RunCommand("./", "ps", "aux")
	//cmd.RunCommand("./", "ls", "-la")
	//cmd.RunCommand("./", "watch", "-n", "2", "ls -la")
}
