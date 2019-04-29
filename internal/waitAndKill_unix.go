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

// +build !windows

package internal

import (
	"fmt"
	"os"
	"time"
)

const sleepDuration = time.Duration(1000) * time.Millisecond

func WaitForPidAndKillProcess(pid int, process *os.Process) {
	for {
		processKilled := checkProcFileOfPidAndKillProcess(pid, process)
		if processKilled {
			return
		}
		time.Sleep(sleepDuration)
	}
}

func checkProcFileOfPidAndKillProcess(pid int, process *os.Process) (processKilled bool) {
	f, err := os.OpenFile(fmt.Sprintf("/proc/%d/status", pid), os.O_RDONLY, 0444)
	defer func() { _ = f.Close() }()

	if err != nil {
		_ = process.Kill()
		return true
	}
	return false
}
