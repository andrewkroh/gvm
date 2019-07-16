package gvm

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

type command struct {
	Path   string
	Args   []string
	Dir    string
	Env    []string
	Stdout func(string)
	Stderr func(string)
}

func makeCommand(cmd string, args ...string) *command {
	return &command{
		Path: cmd,
		Args: args,
	}
}

func (c *command) WithLogger(log logrus.FieldLogger) *command {
	if c.Stdout == nil {
		c.Stdout = infoOutLog(log)
	}
	if c.Stderr == nil {
		c.Stderr = errOutLog(log)
	}
	return c
}

func (c *command) WithDir(dir string) *command {
	c.Dir = dir
	return c
}

func infoOutLog(log logrus.FieldLogger) func(string) {
	return makeOutLog(log.Info)
}

func errOutLog(log logrus.FieldLogger) func(string) {
	return makeOutLog(log.Error)
}

func makeOutLog(fn func(...interface{})) func(string) {
	return func(text string) { fn(text) }
}

func (c *command) Exec() error {
	cmd := exec.Command(c.Path, c.Args...)
	cmd.Dir = c.Dir

	if len(c.Env) > 0 {
		cmd.Env = append(os.Environ(), c.Env...)
	}

	var err error
	var stdout, stderr io.ReadCloser
	if c.Stdout != nil {
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			return err
		}
		defer stdout.Close()
	}

	if c.Stderr != nil {
		stderr, err = cmd.StderrPipe()
		if err != nil {
			return err
		}
		defer stderr.Close()
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	captureLines := func(in io.Reader, fn func(string)) {
		defer wg.Done()
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			text := scanner.Text()
			fn(text)
		}
	}

	if stdout != nil {
		wg.Add(1)
		go captureLines(stdout, c.Stdout)
	}
	if stderr != nil {
		wg.Add(1)
		go captureLines(stderr, c.Stderr)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
