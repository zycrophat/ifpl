/*
 MIT License

 Copyright (c) 2019 Andreas Steffan

 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 SOFTWARE.
*/

package internal

import (
	"log"
	"os"
)

type KillFunc = func(*os.Process) error

func Kill(process *os.Process) error {
	return process.Kill()
}

func GetSendSignalFunc(signal os.Signal) KillFunc {
	return func(process *os.Process) error {
		return process.Signal(signal)
	}
}

func GetLoggingKillFunc(kill KillFunc) KillFunc {
	return func(process *os.Process) error {
		log.Printf("Killing process %d", process.Pid)
		err := kill(process)
		if err != nil {
			log.Printf("Error encountered when killing process %d\n", process.Pid)
		}
		return err
	}
}
