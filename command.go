package command

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type Command struct {
	StdOutWriter  func([]byte)
	StdErrWriter  func([]byte)
	StdInData     string
	Sleep         time.Duration
	OutputAndQuit bool
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
	if t.StdInData != "" {
		in, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer in.Close()
			in.Write([]byte(t.StdInData))
		}()
	}

	c := make(chan bool)

	go func() {
		for {
			select {
			case <-c:
				break
			default:
				// Çıktıları yakalayın ve bir işleve gönderin
				output := bout.Bytes()
				if t.StdOutWriter != nil && len(output) > 0 {
					t.StdOutWriter(output)
					if t.OutputAndQuit {
						break
					}
				}
				if t.Sleep > 0 {
					time.Sleep(t.Sleep * time.Millisecond)
				}
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		c <- true
		log.Println(err)
	}
}
