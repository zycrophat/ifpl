ifpl -- if process lives
========================

Description
-----------

ifpl runs the given command (say, as process p1) and waits for a process with a specified pid (process p2) to terminate.
Upon termination of p2, ifpl will kill p1.

Usage
-----

`$ ifpl [-help] [-pid PID] [-s SIGNAL] CMD [ARGS ...]`

Where:

  - `-help` displays a help message
  - `-pid PID` specifies the PID of the process to wait for to terminate. Defaults to ppid of ifpl.
  - `-s SIGNAL` specifies a signal to send to the `CMD` child process
  - `CMD [ARGS ...]` specifies the CMD and optionally arguments
  
ifpl runs the given `CMD` and waits for the process with pid `PID` to terminate.
Upon termination, the `CMD` child process will be killed.
When the `-s` option is specified with a non-negative value, the given signal will be sent to the CMD child process.
Otherwise go's [os.Process.Kill()](https://golang.org/pkg/os/#Process.Kill) will be used.


Compatibility
-------------

Currently, ifpl runs on Microsoft Windows.
Support for Unix/Linux has been added, but is experimental.


How to build
------------

Prerequisites:

- Go
- Make (optional)

To compile:

```
$ cd cmd/ifpl
$ go build
```

To build distributable archives for Linux & Windows:
```
$ make dist-linux dist-windows
```
