package goparallel

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/kohkimakimoto/goparallel/goparallel/ltsv"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var (
	ErrorTimeout           = errors.New("operation timeout")
	ErrorCmdUndefined      = errors.New("cmd undefined.")
	ErrorUnsupportedFormat = errors.New("unsupported format.")
)

type job struct {
	Cmd    string `yaml:"cmd",json:"cmd"`
	Prefix string `yaml:"prefix",json:"prefix"`
}

func Start() error {
	var timeoutFlag int64
	var versionFlag bool
	var jobsFormat string

	flag.BoolVar(&versionFlag, "v", false, "")
	flag.BoolVar(&versionFlag, "version", false, "")
	flag.StringVar(&jobsFormat, "f", "ltsv", "")
	flag.StringVar(&jobsFormat, "format", "ltsv", "")
	flag.Int64Var(&timeoutFlag, "t", 0, "")
	flag.Int64Var(&timeoutFlag, "timeout", 0, "")

	flag.Usage = printUsage
	flag.Parse()

	if versionFlag {
		fmt.Printf("%s (%s)\n", Version, CommitHash)
		return nil
	}

	var reader io.Reader

	if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		// the stdin is pipe
		reader = os.Stdin
	} else {
		if len(flag.Args()) == 0 {
			flag.Usage()
			return nil
		}
		reader = bytes.NewBufferString(flag.Args()[0])
	}

	var jobs []*job
	var err error
	if jobsFormat == "ltsv" {
		jobs, err = loadJobsFromLTSV(reader)
		if err != nil {
			return err
		}
	} else if jobsFormat == "yaml" {
		jobs, err = loadJobsFromYAML(reader)
		if err != nil {
			return err
		}
	} else if jobsFormat == "json" {
		jobs, err = loadJobsFromJSON(reader)
		if err != nil {
			return err
		}
	} else {
		return ErrorUnsupportedFormat
	}

	wg := &sync.WaitGroup{}

	errchan := make(chan error)

	for _, j := range jobs {
		wg.Add(1)
		go func(j *job) {
			var shell, flag string
			if runtime.GOOS == "windows" {
				shell = "cmd"
				flag = "/C"
			} else {
				shell = "/bin/sh"
				flag = "-c"
			}

			command := j.Cmd

			cmd := exec.Command(shell, flag, command)
			cmd.Stdin = nil

			if j.Prefix != "" {
				cmd.Stdout = newPrefixWriter(os.Stdout, j.Prefix+" ")
				cmd.Stderr = newPrefixWriter(os.Stderr, j.Prefix+" ")
			} else {
				cmd.Stdout = newPrefixWriter(os.Stdout, "")
				cmd.Stderr = newPrefixWriter(os.Stderr, "")
			}

			if err := cmd.Start(); err != nil {
				errchan <- err
				return
			}

			addProcess(cmd.Process)
			defer deleteProecess(cmd.Process)

			if err := cmd.Wait(); err != nil {
				errchan <- err
				return
			}

			wg.Done()
		}(j)
	}

	wait := make(chan bool)
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		wait <- true
	}(wg)

	timeout := make(chan bool)
	if timeoutFlag != 0 {
		go func() {
			time.Sleep(time.Duration(timeoutFlag) * time.Second)
			timeout <- true
		}()
	}

	select {
	case err := <-errchan:
		return err
	case <-wait:
		return nil
	case <-timeout:
		stopChildren()
		return ErrorTimeout
	}

	return nil
}

func loadJobsFromLTSV(r io.Reader) ([]*job, error) {
	lr := ltsv.NewReader(r)
	records, err := lr.ReadAll()
	if err != nil {
		return nil, err
	}

	jobs := []*job{}
	for _, record := range records {
		j := &job{}
		for k, v := range record {
			if k == "cmd" {
				j.Cmd = v
			}

			if k == "prefix" {
				j.Prefix = v
			}
		}

		if j.Cmd == "" {
			return nil, ErrorCmdUndefined
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func loadJobsFromYAML(r io.Reader) ([]*job, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	jobs := []*job{}
	err := yaml.Unmarshal(buf.Bytes(), &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func loadJobsFromJSON(r io.Reader) ([]*job, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	jobs := []*job{}
	err := json.Unmarshal(buf.Bytes(), &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

var processes = map[*os.Process]*os.Process{}
var processesMux = new(sync.Mutex)

func addProcess(proc *os.Process) {
	processesMux.Lock()
	defer processesMux.Unlock()

	processes[proc] = proc
}

func deleteProecess(proc *os.Process) {
	processesMux.Lock()
	defer processesMux.Unlock()

	delete(processes, proc)
}

func stopChildren() {

	// TODO: fix a bug
	//   does not work correctly.
	//   I use "sh -c" to run a command.
	//   The process that works actual task and I need to stop is "grandchild" process.
	//   At now implementation does not stop "grandchild" process. stops "child" process.

	fmt.Fprintf(os.Stderr, "caught a timeout! trying to stop child processes...\n")
	for _, proc := range processes {
		fmt.Fprintf(os.Stderr, "try to kill prosess (pid:%d)\n", proc.Pid)
		err := proc.Kill()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
	}
}

func printUsage() {
	fmt.Println(`Usage: goparallel [<options>] [<commands>]

  goparallel -- Execute commands in parallel.
  version ` + Version + ` (` + CommitHash + `)

  Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
  The MIT License (MIT)
  https://github.com/kohkimakimoto/goparallel

Options:
  -v, --version                 Print version.
  -h, --help                    Show help.
  -f, --format=<jobs format>    Jobs format (ltsv|yaml|json). default 'ltsv'.
  -t, --timeout=<sec>           Timeout seconds.

Commands:
  The commands argument (or stdin stream) is a list of strings include commands that are executed in parallel.
  This is defined in a LTSV format at default. Please see the README that is in the following URL.

    https://github.com/kohkimakimoto/goparallel
`)
}
