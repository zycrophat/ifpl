[![Build Status](https://dev.azure.com/zycrophat/ifpl/_apis/build/status/zycrophat.ifpl?branchName=master)](https://dev.azure.com/zycrophat/ifpl/_build)
[![Release Status](https://vsrm.dev.azure.com/zycrophat/_apis/public/Release/badge/96fe8055-2206-46dc-8be0-0418979b43cd/1/1)](https://dev.azure.com/zycrophat/ifpl/_release?definitionId=1)

ifpl -- if process lives
========================

Description
-----------

__ifpl__ runs the given command (say, as process p1) and waits for a process with a specified pid (process p2) to terminate.
Upon termination of p2, __ifpl__ will kill p1.

Usage
-----

`$ ifpl [-h] [-p PID] [-s SIGNAL] [-v] [-l LOGFILE] CMD [ARGS ...]`

Where:

  - `-h` displays a help message
  - `-p PID` specifies the PID of the process to wait for to terminate. Defaults to ppid of ifpl.
  - `-s SIGNAL` specifies a signal to send to the `CMD` child process
  - `CMD [ARGS ...]` specifies the CMD and optionally arguments
  - `-v` prints ifpl log messages to stdout (unless `-l` is set)
  - `-l LOGFILE` file to write log messages to (implies `-v`)

__ifpl__ runs the given `CMD` and waits for the process with pid `PID` to terminate.
Upon termination, the `CMD` child process will be killed.
When the `-s` option is specified with a non-negative value, the given signal will be sent to the CMD child process.
Otherwise go's [os.Process.Kill()](https://golang.org/pkg/os/#Process.Kill) will be used.

__ifpl__ returns with the exit status of `CMD` or with 255 if an error occurred.

Compatibility
-------------

Currently, __ifpl__ runs on Microsoft Windows.
For Unix/Linux procfs is required.


How to build
------------

Prerequisites:

- Go1.12
- Make (optional)

To compile:

```
$ cd cmd/ifpl
$ go build
```

To build distributable archives for Linux & Windows:
```
$ make dist-windows-x86-64 dist-linux-x86-64 dist-linux-arm64
```
