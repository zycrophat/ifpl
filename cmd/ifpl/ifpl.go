/*
 * MIT License
 *
 * Copyright (c) 2019 Andreas Steffan
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"flag"
	"fmt"
	"ifpl/internal"
	"os"
	"os/exec"
	"os/signal"
)

const (
	exitCodeOffset = 240
)

const (
	pidFlagName  = "pid"
	helpFlagName = "help"
)

func main() {
	ifplArgs := parseArgs()
	if ifplArgs.help {
		printHelp()
		os.Exit(0)
	}

	if ifplArgs.cmdName == "" {
		printHelp()
		os.Exit(exitCodeOffset + 1)
	}

	cmd := createAndConfigureCmd(ifplArgs.cmdName, ifplArgs.cmdArgs)

	err := cmd.Start()
	if err != nil {
		os.Exit(exitCodeOffset + 2)
	}

	go redirectSignals(cmd.Process)
	go internal.WaitForPidAndKillProcess(ifplArgs.pid, cmd.Process)

	_ = cmd.Wait()

	os.Exit(cmd.ProcessState.ExitCode())
}

func printHelp() {
	fmt.Print("ifpl -- if process lives\n")
	fmt.Printf("(pid: %d, ppid %d)\n\n", os.Getpid(), os.Getppid())

	fmt.Print("Usage:\n")
	fmt.Printf("ifpl [-%s] [-%s <pid>] CMD [ARGS ...]\n", helpFlagName, pidFlagName)
	flag.PrintDefaults()

	fmt.Print("\nifpl runs the given CMD and waits for the process with pid <pid> to terminate. " +
		"Upon termination CMD will be killed.\n\n")

}

type ifplArgs struct {
	pid     int      // the pid to wait for to terminate
	cmdName string   // the cmd to execute
	cmdArgs []string // args for the cmd to execute

	help bool // flag to print help
}

func parseArgs() ifplArgs {
	pid := flag.Int(pidFlagName, os.Getppid(), "the pid to wait for to terminate. "+
		"Defaults to ppid of ifpl")
	helpIsRequested := flag.Bool(helpFlagName, false, "displays this help message")
	flag.Parse()
	args := flag.Args()

	ifplArgs := ifplArgs{
		pid: *pid,
		cmdName: func() string {
			if len(args) > 0 {
				return args[0]
			} else {
				return ""
			}
		}(),
		cmdArgs: func() []string {
			if len(args) > 1 {
				return args[1:]
			} else {
				return []string{}
			}
		}(),
		help: *helpIsRequested,
	}
	return ifplArgs
}

func createAndConfigureCmd(cmdName string, cmdArgs []string) *exec.Cmd {
	cmd := exec.Command(cmdName, cmdArgs...)
	configureCmd(cmd)
	return cmd
}

func configureCmd(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func redirectSignals(process *os.Process) {
	c := make(chan os.Signal)
	signal.Notify(c)

	for {
		s := <-c
		_ = process.Signal(s)
	}
}