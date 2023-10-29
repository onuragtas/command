package command

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Command struct {
	stdInFunction func()
	stdInDuration int

	stdOutWriter func([]byte)
	stdErrWriter func([]byte)
	stdInData    string
}

func (t *Command) RunCommand(path string, name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	fmt.Println("command:", name, arg)
	if path != "" {
		cmd.Dir = path
	}

	var bout, berr bytes.Buffer
	cmd.Stdout = &bout
	cmd.Stderr = &berr

	// Komutun stdin'e yazılması
	if t.stdInData != "" {
		in, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer in.Close()
			in.Write([]byte(t.stdInData))
		}()
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			// Çıktıları yakalayın ve bir işleve gönderin
			output := bout.Bytes()
			if t.stdOutWriter != nil {
				t.stdOutWriter(output)
			}

			// Çıktıları sıfırlamayın, bu şekilde sürekli olarak çıktıları izleyebilirsiniz

			// Belirli bir aralıkta veya her çıktı sonrasında bir süre bekleme ekleyebilirsiniz
			// time.Sleep(1 * time.Second)
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}
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
