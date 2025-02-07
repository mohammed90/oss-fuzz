// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CopyFile(src string, dst string) {
	contents, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(dst, contents, 0644)
	if err != nil {
		panic(err)
	}
}

func TryFixCCompilation(cmdline []string) bool {
	var newFile string = ""
	for i, arg := range cmdline {
		if !strings.HasSuffix(arg, ".c") {
			continue
		}
		if _, err := os.Stat(arg); errors.Is(err, os.ErrNotExist) {
			continue
		}
		newFile = strings.TrimSuffix(arg, ".c")
		newFile += ".cpp"
		CopyFile(arg, newFile)
		cmdline[i] = newFile
		break
	}
	if newFile == "" {
		return false
	}
	cmd := exec.Command("clang++", cmdline...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	fmt.Println(cmd)
	err := cmd.Run()
	fmt.Println(outb.String())
	fmt.Println(errb.String())
	if err != nil {
		os.Exit(cmd.ProcessState.ExitCode())
	}
	return true
}

func main() {
	args := os.Args[1:]
	basename := filepath.Base(os.Args[0])
	isCPP := basename == "clang++"
	newArgs := []string{"-w"}
	newArgs = append(args, newArgs...)
	var cmd *exec.Cmd
	if isCPP {
		cmd = exec.Command("clang++", newArgs...)
	} else {
		cmd = exec.Command("clang", newArgs...)
	}
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()

	if err == nil {
		os.Exit(0)
	}

	if isCPP || !TryFixCCompilation(newArgs) {
		// Nothing else we can do. Just print the error and exit.
		fmt.Println(outb.String())
		fmt.Println(errb.String())
		os.Exit(cmd.ProcessState.ExitCode())
	}
}
