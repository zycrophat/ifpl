ifpl -- if process lives
========================

Description
-----------

ifpl runs the given command and waits for the process with pid \<pid> to terminate.
Upon termination of the process with pid \<pid> the spawned process will be killed.

Usage
-----

`ifpl.exe [-help] [-pid <pid>] CMD [ARGS ...]`

Where:

  - `-help` displays a help message
  - `-pid <pid>` specifies the pid of the process to wait for to terminate. Defaults to ppid of ifpl.
  - `CMD [ARGS ...]` specifies the CMD and optionally arguments

Compatibility
-------------

Currently, ifpl runs on Microsoft Windows.
Support for Unix/Linux has been added, but is experimental.
