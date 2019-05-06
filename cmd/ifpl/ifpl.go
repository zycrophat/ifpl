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
	"syscall"
)

const (
	exitCodeOffset = 240
)

const (
	pidFlagName       = "pid"
	helpFlagName      = "help"
	signalFlagName    = "s"
	signalFlagDefault = -1
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

	killFunc := getKillFunc(ifplArgs.signal)

	go redirectSignals(cmd.Process)
	go internal.WaitForPidAndKillProcess(ifplArgs.pid, cmd.Process, killFunc)

	_ = cmd.Wait()

	os.Exit(cmd.ProcessState.ExitCode())
}

func printHelp() {
	fmt.Print("ifpl -- if process lives\n")
	fmt.Printf("(pid: %d, ppid %d)\n\n", os.Getpid(), os.Getppid())

	fmt.Print("Usage:\n")
	fmt.Printf("ifpl [-%s] [-%s PID] [-%s SIGNAL] CMD [ARGS ...]\n", helpFlagName, pidFlagName, signalFlagName)
	flag.PrintDefaults()

	fmt.Print("\nifpl runs the given CMD and waits for the process with pid PID to terminate.\n")
	fmt.Print("Upon termination, the CMD child process will be killed.\n")
	fmt.Printf("When the -%s option is specified with a non-negative value, the given signal will be sent to the CMD child process.\n", signalFlagName)
	fmt.Print("Otherwise go's os.Process.Kill() will be used.\n\n")

}

type ifplArgs struct {
	pid     int      // the pid to wait for to terminate
	cmdName string   // the cmd to execute
	cmdArgs []string // args for the cmd to execute
	signal  int

	help bool // flag to print help
}

func parseArgs() ifplArgs {
	pid := flag.Int(pidFlagName, os.Getppid(), "the pid to wait for to terminate. "+
		"Defaults to ppid of ifpl")
	help := flag.Bool(helpFlagName, false, "displays this help message")
	signalArg := flag.Int(signalFlagName, signalFlagDefault, "signal to be sent to CMD")

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
		signal: *signalArg,
		help:   *help,
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
		go process.Signal(s)
	}
}

func getKillFunc(signalArg int) internal.KillFunc {
	var kill internal.KillFunc
	if signalArg < 0 {
		kill = internal.Kill
	} else {
		kill = internal.GetSendSignalFunc(syscall.Signal(signalArg))
	}

	return kill
}
