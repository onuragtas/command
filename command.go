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

func (c *Command) Run(cmd string) ([]byte, error) {
	return exec.Command("bash", cmd).Output()
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
	cError := make(chan bool)

	if err := cmd.Start(); err != nil {
		c <- true
		cError <- true
	}

	go func() {
		loop := true
		for loop {
			select {
			case <-c:
				loop = false
				break
			default:
				// Çıktıları yakalayın ve bir işleve gönderin
				loop = t.send(bout, t.StdOutWriter)
			}
		}
	}()

	go func() {
		loop := true
		for loop {
			select {
			case <-cError:
				loop = false
				break
			default:
				// Çıktıları yakalayın ve bir işleve gönderin
				loop = t.send(berr, t.StdErrWriter)
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		c <- true
		cError <- true
	}

	c <- true
	cError <- true

	t.send(bout, t.StdOutWriter)
	t.send(berr, t.StdErrWriter)
}

func (t *Command) send(buf bytes.Buffer, writer func([]byte)) bool {
	loop := true
	output := buf.Bytes()
	if writer != nil && len(output) > 0 {
		writer(output)
		if t.OutputAndQuit {
			loop = false
		}
	}
	if t.Sleep > 0 {
		time.Sleep(t.Sleep * time.Millisecond)
	}

	return loop
}
