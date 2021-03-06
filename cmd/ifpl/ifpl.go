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
	"github.com/zycrophat/ifpl/internal"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

const (
	ifplErrorExitCode = 255
)

const (
	pidFlagName         = "p"
	helpFlagName        = "h"
	signalFlagName      = "s"
	signalFlagDefault   = -1
	verboseFlagName     = "v"
	logFilePathFlagName = "l"
)

func main() {
	ifplArgs := parseArgs()

	f := configureLog(ifplArgs)
	if f != nil {
		defer func() { _ = f.Close() }()
	}

	if ifplArgs.help {
		printHelp()
		os.Exit(0)
	}

	if ifplArgs.cmdName == "" {
		printHelp()
		os.Exit(ifplErrorExitCode)
	}

	killFunc := getKillFunc(ifplArgs.signal)

	exitCode := startAndWaitForCmd(ifplArgs, killFunc)

	os.Exit(exitCode)
}

func startAndWaitForCmd(args ifplArgs, killFunc internal.KillFunc) int {
	cmd := createAndConfigureCmd(args.cmdName, args.cmdArgs)

	log.Printf("Starting command: %s %s\n", cmd.Path, strings.Join(cmd.Args[1:], " "))
	err := cmd.Start()
	if err != nil {
		log.Printf("Cannot start command: %s\n", err)
		os.Exit(ifplErrorExitCode)
	}

	go redirectSignals(cmd.Process)
	go internal.WaitForPidAndKillProcess(args.pid, cmd.Process, killFunc)
	log.Printf("Waiting for process %d to terminate\n", args.pid)

	log.Printf("Waiting for command to terminate\n")
	_ = cmd.Wait()
	log.Printf("Command has terminated\n")

	return cmd.ProcessState.ExitCode()
}

func printHelp() {
	fmt.Print("ifpl ─ if process lives\n")
	fmt.Printf("(pid: %d, ppid %d)\n\n", os.Getpid(), os.Getppid())

	fmt.Print("Usage:\n")
	fmt.Printf("ifpl [-%s] [-%s PID] [-%s SIGNAL] [-%s] [-%s LOGFILE] CMD [ARGS ...]\n",
		helpFlagName, pidFlagName, signalFlagName, verboseFlagName, logFilePathFlagName)
	flag.PrintDefaults()
}

type ifplArgs struct {
	pid        int      // the pid to wait for to terminate
	cmdName    string   // the cmd to execute
	cmdArgs    []string // args for the cmd to execute
	signal     int
	verboseArg bool

	help           bool // flag to print help
	isWriteLogFile bool
	logFilePathArg string
}

func parseArgs() ifplArgs {
	pid := flag.Int(pidFlagName, os.Getppid(), "the pid to wait for to terminate. "+
		"Defaults to ppid of ifpl")
	help := flag.Bool(helpFlagName, false, "displays this help message")
	signalArg := flag.Int(signalFlagName, signalFlagDefault, "signal to be sent to CMD")
	verboseArg := flag.Bool(verboseFlagName, false, "prints ifpl log messages to stdout")
	logFilePathArg := flag.String(logFilePathFlagName, "", "file to write log messages to")

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
		signal:         *signalArg,
		help:           *help,
		verboseArg:     *verboseArg,
		isWriteLogFile: isFlagSet(logFilePathFlagName),
		logFilePathArg: *logFilePathArg,
	}
	return ifplArgs
}

func isFlagSet(flagName string) bool {
	isFlagSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			isFlagSet = true
		}
	})
	return isFlagSet
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
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	for {
		s := <-c
		go func() { _ = process.Signal(s) }()
	}
}

func getKillFunc(signalArg int) internal.KillFunc {
	var kill internal.KillFunc
	if signalArg < 0 {
		kill = internal.Kill
	} else {
		kill = internal.GetSendSignalFunc(syscall.Signal(signalArg))
	}

	killAndLog := internal.GetLoggingKillFunc(kill)
	return killAndLog
}

func configureLog(args ifplArgs) io.Closer {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	consoleWriter := func() io.Writer {
		if args.verboseArg {
			return os.Stdout
		} else {
			return ioutil.Discard
		}
	}()

	if args.isWriteLogFile {
		f, err := os.OpenFile(args.logFilePathArg, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Printf("Cannot write log to file: %s\n", err)
			os.Exit(ifplErrorExitCode)
		}
		multiWriter := io.MultiWriter(f, consoleWriter)
		log.SetOutput(multiWriter)

		return f
	} else {
		log.SetOutput(consoleWriter)
	}

	return nil
}
