ifpl -- if process lives
========================

Description
-----------

ifpl runs the given command (say, as process p1) and waits for a process with a specified pid (process p2) to terminate.
Upon termination of p2, ifpl will kill p1.

Usage
-----

`$ ifpl [-help] [-pid PID] CMD [ARGS ...]`

Where:

  - `-help` displays a help message
  - `-pid PID` specifies the PID of the process to wait for to terminate. Defaults to ppid of ifpl.
  - `CMD [ARGS ...]` specifies the CMD and optionally arguments

Compatibility
-------------

Currently, ifpl runs on Microsoft Windows.
Support for Unix/Linux has been added, but is experimental.
